package internal

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/chanxuehong/log"
	"github.com/chanxuehong/mns.aliyun.v20150606"
)

func DoHTTP(ctx context.Context, httpMethod string, _url *url.URL, header http.Header, reqBody []byte, respBuffer *bytes.Buffer, config mns.Config) (requestId string, statusCode int, respBody []byte, err error) {
	logger, _ := log.FromContext(ctx)
	for i := 0; i < 3; i++ {
		respBuffer.Reset()
		requestId, statusCode, respBody, err = doHTTP(ctx, httpMethod, _url, header, reqBody, respBuffer, config)
		if err == nil {
			return
		}
		if logger != nil {
			logger.Error("mns: DoHTTP encountered an error", "error-type", reflect.TypeOf(err).String(), "error", err.Error())
		}
		if !shouldRetryRequest(err) {
			return
		}
	}
	return
}

// shouldRetryRequest 根据 err 判断是否需要重试
//
// 由于双方的 keepalive 等参数配置不一样, 可能服务器端会关闭一些 connection 而客户端没有及时发现, 会抛出一些错误, 一般通过重试可以正常的工作
//
// 随着标准库的更新可能会变化
func shouldRetryRequest(err error) bool {
	if _, ok := err.(*net.OpError); ok {
		if strings.Contains(err.Error(), "connection reset by peer") {
			return true
		}
		return false
	}
	if uerr, ok := err.(*url.Error); ok {
		if uerr.Err == io.EOF {
			return true
		}
		return false
	}
	return false
}

func doHTTP(ctx context.Context, httpMethod string, _url *url.URL, header http.Header, reqBody []byte, respBuffer *bytes.Buffer, config mns.Config) (requestId string, statusCode int, respBody []byte, err error) {
	if httpMethod == "" {
		httpMethod = http.MethodGet
	}
	if header == nil {
		header = make(http.Header, 8)
	}
	if respBuffer == nil {
		respBuffer = bytes.NewBuffer(make([]byte, 0, 16<<10))
	}

	header.Set("Date", FormatDate(time.Now()))
	header.Set("X-Mns-Version", Version)
	header.Set("Content-Type", ContentType)
	if len(reqBody) > 0 {
		header.Set("Content-Md5", ContentMD5(reqBody))
	}
	header.Set("Authorization", Authorization(config.AccessKeyId, Sign(httpMethod, header, _url.RequestURI(), config.AccessKeySecret)))

	req := &http.Request{
		Method:        httpMethod,
		URL:           _url,
		Header:        header,
		Host:          _url.Host,
		ContentLength: int64(len(reqBody)),
	}
	if req.ContentLength > 0 {
		req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(reqBody)), nil
		}
	}

	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}
	if ctx != context.Background() {
		req = req.WithContext(ctx)
	}
	resp, err := config.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if _, err = respBuffer.ReadFrom(resp.Body); err != nil {
		return
	}
	requestId = resp.Header.Get("X-Mns-Request-Id")
	return requestId, resp.StatusCode, respBuffer.Bytes(), nil
}
