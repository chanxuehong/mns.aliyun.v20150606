package topic

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chanxuehong/mns.aliyun.v20150606"
	"github.com/chanxuehong/mns.aliyun.v20150606/internal"
)

type Topic struct {
	config mns.Config
	topic  string // http://$AccountId.mns.<Region>.aliyuncs.com/topics/$TopicName
}

// New 创建一个新的 Topic
//  endpoint: http://$AccountId.mns.<Region>.aliyuncs.com
//  topic:    topic name
func New(endpoint, topic string, config mns.Config) *Topic {
	endpoint = strings.TrimRight(endpoint, "/")
	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}
	return &Topic{
		config: config,
		topic:  endpoint + "/topics/" + topic,
	}
}

type PublishMessageRequest struct {
	XMLName struct{} `xml:"Message"`

	MessageBody       []byte      `xml:"MessageBody"`
	MessageTag        string      `xml:"MessageTag,omitempty"`
	MessageAttributes interface{} `xml:"MessageAttributes,omitempty"`
}

type PublishMessageResponse struct {
	XMLName struct{} `xml:"Message"`

	MessageId      string `xml:"MessageId"`
	MessageBodyMD5 string `xml:"MessageBodyMD5"`
}

func (t *Topic) PublishMessage(msg *PublishMessageRequest) (requestId string, resp *PublishMessageResponse, err error) {
	return t.PublishMessageContext(context.Background(), msg)
}

func (t *Topic) PublishMessageContext(ctx context.Context, msg *PublishMessageRequest) (requestId string, resp *PublishMessageResponse, err error) {
	if msg == nil || len(msg.MessageBody) == 0 {
		err = errors.New("the MessageBody must not be empty")
		return
	}
	if len(msg.MessageTag) > 16 {
		err = errors.New("the length of MessageTag cannot be greater than 16")
		return
	}
	if t.config.Base64Enabled {
		msg.MessageBody = internal.Base64Encode(msg.MessageBody)
	}

	_url, err := internal.ParseURL(t.topic + "/messages")
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	reqBuffer := pool.Get()
	defer pool.Put(reqBuffer)
	reqBuffer.Reset()
	if err = xml.NewEncoder(reqBuffer).Encode(msg); err != nil {
		return
	}
	reqBody := reqBuffer.Bytes()

	pool = mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodPost, _url, nil, reqBody, respBuffer, t.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result PublishMessageResponse
		if err = xml.Unmarshal(respBody, &result); err != nil {
			return
		}
		if want := internal.MessageBodyMD5(msg.MessageBody); strings.ToUpper(result.MessageBodyMD5) != want {
			err = fmt.Errorf("the MessageBodyMD5 mismatch, have:%s, want:%s", result.MessageBodyMD5, want)
			return
		}
		resp = &result
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}
