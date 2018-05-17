package internal

import (
	"encoding/base64"
	"reflect"
)

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
	return "cannot unmarshal xml " + string(e.Data) + " into Go value of type " + e.Type.String() + ", " + e.Err.Error()
}

func NewMessageBodyMD5MismatchError(messageBody []byte, have, want string) error {
	return &MessageBodyMD5MismatchError{
		MessageBody: messageBody,
		Have:        have,
		Want:        want,
	}
}

type MessageBodyMD5MismatchError struct {
	MessageBody []byte
	Have        string
	Want        string
}

func (e *MessageBodyMD5MismatchError) Error() string {
	return "the MessageBodyMD5 mismatch, have: " + e.Have + ", want: " + e.Want + ", base64(MessageBodyMD5): " + base64.StdEncoding.EncodeToString(e.MessageBody)
}
