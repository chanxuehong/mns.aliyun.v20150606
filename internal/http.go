package internal

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

func DoHTTP(ctx context.Context, httpMethod string, _url *url.URL, header http.Header, reqBody []byte, respBuffer *bytes.Buffer, config mns.Config) (requestId string, statusCode int, respBody []byte, err error) {
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
	header.Set("Host", _url.Host)
	header.Set("X-Mns-Version", Version)
	header.Set("Content-Length", strconv.Itoa(len(reqBody)))
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
