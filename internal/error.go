package internal

import (
	"encoding/xml"
	"reflect"

	"github.com/chanxuehong/mns.aliyun.v20150606"
)

// UnmarshalError unmarshal error response.
func UnmarshalError(requestId string, statusCode int, body []byte) error {
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

func NewXMLUnmarshalError(data []byte, v interface{}, err error) error {
	return &XMLUnmarshalError{
		Data: data,
		Type: reflect.TypeOf(v),
		Err:  err,
	}
}

type XMLUnmarshalError struct {
	Data []byte
	Type reflect.Type
	Err  error
}

func (e *XMLUnmarshalError) Error() string {
	return "xml: cannot unmarshal " + string(e.Data) + " into Go value of type " + e.Type.String() + ": " + e.Err.Error()
}
