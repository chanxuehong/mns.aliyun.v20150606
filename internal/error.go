package internal

import (
	"encoding/xml"
	"fmt"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

func UnmarshalError(requestId string, statusCode int, body []byte) error {
	var result mns.Error
	if err := xml.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("unmarshal error response failed, response=%q, error=%s", body, err.Error())
	}
	if result.RequestId == "" {
		result.RequestId = requestId
	}
	result.HttpStatusCode = statusCode
	return &result
}
