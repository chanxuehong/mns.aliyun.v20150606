package mns

import (
	"encoding/xml"
)

const (
	ErrorHttpStatusCodeQueueNotExist      = 404
	ErrorHttpStatusCodeTopicNotExist      = 404
	ErrorHttpStatusCodeMessageNotExist    = 404
	ErrorHttpStatusCodeReceiptHandleError = 400
)

const (
	ErrorCodeQueueNotExist      = "QueueNotExist"
	ErrorCodeTopicNotExist      = "TopicNotExist"
	ErrorCodeMessageNotExist    = "MessageNotExist"
	ErrorCodeReceiptHandleError = "ReceiptHandleError"
)

func IsQueueNotExist(err error) bool {
	v, ok := err.(*Error)
	if !ok {
		return false
	}
	if v == nil {
		return false
	}
	return v.HttpStatusCode == ErrorHttpStatusCodeQueueNotExist && v.Code == ErrorCodeQueueNotExist
}

func IsTopicNotExist(err error) bool {
	v, ok := err.(*Error)
	if !ok {
		return false
	}
	if v == nil {
		return false
	}
	return v.HttpStatusCode == ErrorHttpStatusCodeTopicNotExist && v.Code == ErrorCodeTopicNotExist
}

func IsMessageNotExist(err error) bool {
	v, ok := err.(*Error)
	if !ok {
		return false
	}
	if v == nil {
		return false
	}
	return v.HttpStatusCode == ErrorHttpStatusCodeMessageNotExist && v.Code == ErrorCodeMessageNotExist
}

func IsReceiptHandleError(err error) bool {
	v, ok := err.(*Error)
	if !ok {
		return false
	}
	if v == nil {
		return false
	}
	return v.HttpStatusCode == ErrorHttpStatusCodeReceiptHandleError && v.Code == ErrorCodeReceiptHandleError
}

var _ error = (*Error)(nil)

// Error 表示 MNS 的错误响应.
type Error struct {
	XMLName struct{} `xml:"Error" json:"-"`

	HttpStatusCode int    `xml:"HttpStatusCode" json:"http_status_code"` // HTTP 状态码
	Code           string `xml:"Code" json:"code"`                       // MNS 返回给用户的错误码。
	Message        string `xml:"Message" json:"message"`                 // MNS 给出的详细错误信息。
	RequestId      string `xml:"RequestId" json:"request_id"`            // 用于唯一标识该次请求的编号；当你无法解决问题时，可以提供这个 RequestId 寻求 MNS 支持工程师的帮助。
	HostId         string `xml:"HostId" json:"host_id"`                  // 用于标识访问的 MNS 服务的地域。
}

func (e *Error) Error() string {
	b, _ := xml.Marshal(e)
	return string(b)
}
