package internal

import (
	"encoding/xml"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

// UnmarshalErrorResponse 解析 mns 的标准返回到 *mns.Error.
func UnmarshalErrorResponse(requestId string, statusCode int, body []byte) error {
	var result mns.Error
	if err := xml.Unmarshal(body, &result); err != nil {
		return NewXMLUnmarshalError(body, &result, err)
	}
	if result.RequestId == "" {
		result.RequestId = requestId
	}
	result.HttpStatusCode = statusCode
	return &result
}
