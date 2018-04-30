package internal

import (
	"encoding/xml"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

func UnmarshalError(requestId string, statusCode int, body []byte) error {
	var result mns.Error
	if err := xml.Unmarshal(body, &result); err != nil {
		return err
	}
	if result.RequestId == "" {
		result.RequestId = requestId
	}
	result.HttpStatusCode = statusCode
	return &result
}
