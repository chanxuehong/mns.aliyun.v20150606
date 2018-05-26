package mns

import (
	"net/http"
	"time"

	"github.com/chanxuehong/mns.aliyun.v20150606/log"
)

type Config struct {
	AccessKeyId     string
	AccessKeySecret string

	// following is optional
	Timeout       time.Duration
	Base64Enabled bool
	HttpClient    *http.Client
	Logger        log.Logger
}
