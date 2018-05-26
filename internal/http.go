package internal

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

func DoHTTP(ctx context.Context, httpMethod string, _url *url.URL, header http.Header, reqBody []byte, respBuffer *bytes.Buffer, config mns.Config) (requestId string, statusCode int, respBody []byte, err error) {
	for i := 0; i < 3; i++ {
		respBuffer.Reset()
		requestId, statusCode, respBody, err = doHTTP(ctx, httpMethod, _url, header, reqBody, respBuffer, config)
		switch {
		default:
			return
		case err == nil:
			return
		case shouldRetryRequest(err, config.Logger):
			continue
		}
	}
	return
}

func shouldRetryRequest(err error, lg mns.Logger) bool {
	uerr, ok := err.(*url.Error)
	if !ok {
		return false
	}
	if lg != nil {
		lg.Errorf("DoHTTP encountered an url.Error: %v, cause: %s", err, reflect.TypeOf(uerr.Err).String())
	}
	// TODO: 目前返回的是 io.EOF, 哪天 http 包返回别的错误需要修改这里!!!
	if uerr.Err == io.EOF {
		return true // http 时不时返回 Err == io.EOF 的 url.Error 错误, 一般是阿里云主动关闭了连接导致的, 可以调整参数来解决, 这里也重试几次吧
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
