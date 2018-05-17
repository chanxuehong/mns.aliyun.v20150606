package queue

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chanxuehong/mns.aliyun.v20150606"
	"github.com/chanxuehong/mns.aliyun.v20150606/internal"
)

type Queue struct {
	config mns.Config
	queue  string // http://$AccountId.mns.<Region>.aliyuncs.com/queues/$QueueName
}

// New 创建一个新的 Queue
//  endpoint: http://$AccountId.mns.<Region>.aliyuncs.com
//  queue:    queue name
func New(endpoint, queue string, config mns.Config) *Queue {
	endpoint = strings.TrimRight(endpoint, "/")
	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}
	return &Queue{
		config: config,
		queue:  endpoint + "/queues/" + queue,
	}
}

type SendMessageRequest struct {
	XMLName struct{} `xml:"Message"`

	MessageBody  []byte `xml:"MessageBody"`
	DelaySeconds int    `xml:"DelaySeconds,omitempty"`
	Priority     int    `xml:"Priority,omitempty"`
}

type SendMessageResponse struct {
	XMLName struct{} `xml:"Message"`

	MessageId      string `xml:"MessageId"`
	MessageBodyMD5 string `xml:"MessageBodyMD5"`
	ReceiptHandle  string `xml:"ReceiptHandle,omitempty"`
}

func (q *Queue) SendMessage(msg *SendMessageRequest) (requestId string, resp *SendMessageResponse, err error) {
	return q.SendMessageContext(context.Background(), msg)
}

func (q *Queue) SendMessageContext(ctx context.Context, msg *SendMessageRequest) (requestId string, resp *SendMessageResponse, err error) {
	if msg == nil || len(msg.MessageBody) == 0 {
		err = errors.New("the MessageBody must not be empty")
		return
	}
	if q.config.Base64Enabled {
		msg.MessageBody = internal.Base64Encode(msg.MessageBody)
	}

	_url, err := internal.ParseURL(q.queue + "/messages")
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
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodPost, _url, nil, reqBody, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result SendMessageResponse
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
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

type BatchSendMessageResponseItem struct {
	XMLName struct{} `xml:"Message"`

	ErrorCode    string `xml:"ErrorCode"`
	ErrorMessage string `xml:"ErrorMessage"`

	// 要么上面两个属性有值, 要么下面三个属性有值, 不会同时有值

	MessageId      string `xml:"MessageId"`
	MessageBodyMD5 string `xml:"MessageBodyMD5"`
	ReceiptHandle  string `xml:"ReceiptHandle,omitempty"`
}

func (q *Queue) BatchSendMessage(msgs []SendMessageRequest) (requestId string, resp []BatchSendMessageResponseItem, err error) {
	return q.BatchSendMessageContext(context.Background(), msgs)
}

func (q *Queue) BatchSendMessageContext(ctx context.Context, msgs []SendMessageRequest) (requestId string, resp []BatchSendMessageResponseItem, err error) {
	if len(msgs) < 1 || len(msgs) > 16 {
		err = errors.New("the length of msgs is invalid")
		return
	}
	for i := range msgs {
		if len(msgs[i].MessageBody) == 0 {
			err = errors.New("the MessageBody must not be empty")
			return
		}
	}
	if q.config.Base64Enabled {
		for i := range msgs {
			msgs[i].MessageBody = internal.Base64Encode(msgs[i].MessageBody)
		}
	}

	_url, err := internal.ParseURL(q.queue + "/messages")
	if err != nil {
		return
	}

	var req = struct {
		XMLName struct{} `xml:"Messages"`

		Messages []SendMessageRequest `xml:"Message,omitempty"`
	}{
		Messages: msgs,
	}

	pool := mns.GetBytesBufferPool()
	reqBuffer := pool.Get()
	defer pool.Put(reqBuffer)
	reqBuffer.Reset()
	if err = xml.NewEncoder(reqBuffer).Encode(&req); err != nil {
		return
	}
	reqBody := reqBuffer.Bytes()

	pool = mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodPost, _url, nil, reqBody, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode == 500:
		if !(bytes.Contains(respBody, []byte(`</Messages>`)) && bytes.Contains(respBody, []byte(`</Message>`))) {
			err = internal.UnmarshalError(requestId, statusCode, respBody)
			return
		}
		fallthrough // 只发送了部分消息
	case statusCode/100 == 2:
		var result struct {
			XMLName struct{} `xml:"Messages"`

			Messages []BatchSendMessageResponseItem `xml:"Message"`
		}
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		if len(req.Messages) != len(result.Messages) {
			err = fmt.Errorf("result message count mismatch, have:%d, want:%d", len(result.Messages), len(req.Messages))
			return
		}
		resultMessages := result.Messages
		reqMessages := req.Messages
		reqMessages = reqMessages[:len(resultMessages)]
		for i := range resultMessages {
			if resultMessages[i].ErrorCode != "" {
				continue
			}
			want := internal.MessageBodyMD5(reqMessages[i].MessageBody)
			if strings.ToUpper(resultMessages[i].MessageBodyMD5) != want {
				err = fmt.Errorf("the %d'th MessageBodyMD5 mismatch, have:%s, want:%s", i, resultMessages[i].MessageBodyMD5, want)
				return
			}
		}
		resp = result.Messages
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

type Message struct {
	XMLName struct{} `xml:"Message"`

	MessageId        string `xml:"MessageId"`
	ReceiptHandle    string `xml:"ReceiptHandle"`
	MessageBody      []byte `xml:"MessageBody"`
	MessageBodyMD5   string `xml:"MessageBodyMD5"`
	EnqueueTime      int64  `xml:"EnqueueTime"`
	NextVisibleTime  int64  `xml:"NextVisibleTime"`
	FirstDequeueTime int64  `xml:"FirstDequeueTime"`
	DequeueCount     int    `xml:"DequeueCount"`
	Priority         int    `xml:"Priority"`
}

func (q *Queue) ReceiveMessage(waitSeconds int) (requestId string, msg *Message, err error) {
	return q.ReceiveMessageContext(context.Background(), waitSeconds)
}

func (q *Queue) ReceiveMessageContext(ctx context.Context, waitSeconds int) (requestId string, msg *Message, err error) {
	if waitSeconds < 0 || waitSeconds > 30 {
		waitSeconds = 30
	}
	config := q.config
	if config.Timeout > 0 {
		var timeout time.Duration
		if waitSeconds == 0 {
			timeout = (30 + 10) * time.Second
		} else {
			timeout = time.Duration(waitSeconds+10) * time.Second
		}
		if config.Timeout < timeout {
			config.Timeout = timeout
		}
	}

	rawurl := q.queue + "/messages"
	if waitSeconds > 0 {
		rawurl += "?waitseconds=" + strconv.Itoa(waitSeconds)
	}
	_url, err := internal.ParseURL(rawurl)
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodGet, _url, nil, nil, respBuffer, config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result Message
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		if want := internal.MessageBodyMD5(result.MessageBody); strings.ToUpper(result.MessageBodyMD5) != want {
			err = fmt.Errorf("the MessageBodyMD5 mismatch, have:%s, want:%s", result.MessageBodyMD5, want)
			return
		}
		if q.config.Base64Enabled && len(result.MessageBody) > 0 {
			result.MessageBody, err = internal.Base64Decode(result.MessageBody)
			if err != nil {
				return
			}
			result.MessageBodyMD5 = internal.MessageBodyMD5(result.MessageBody)
		}
		msg = &result
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

func (q *Queue) BatchReceiveMessage(numOfMessages, waitSeconds int) (requestId string, msgs []Message, err error) {
	return q.BatchReceiveMessageContext(context.Background(), numOfMessages, waitSeconds)
}

func (q *Queue) BatchReceiveMessageContext(ctx context.Context, numOfMessages, waitSeconds int) (requestId string, msgs []Message, err error) {
	if numOfMessages < 1 || numOfMessages > 16 {
		numOfMessages = 16
	}
	if waitSeconds < 0 || waitSeconds > 30 {
		waitSeconds = 30
	}
	config := q.config
	if config.Timeout > 0 {
		var timeout time.Duration
		if waitSeconds == 0 {
			timeout = (30 + 10) * time.Second
		} else {
			timeout = time.Duration(waitSeconds+10) * time.Second
		}
		if config.Timeout < timeout {
			config.Timeout = timeout
		}
	}

	rawurl := q.queue + "/messages?numOfMessages=" + strconv.Itoa(numOfMessages)
	if waitSeconds > 0 {
		rawurl += "&waitseconds=" + strconv.Itoa(waitSeconds)
	}
	_url, err := internal.ParseURL(rawurl)
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodGet, _url, nil, nil, respBuffer, config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result struct {
			XMLName struct{} `xml:"Messages"`

			Messages []Message `xml:"Message"`
		}
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		resultMessages := result.Messages
		for i := range resultMessages {
			want := internal.MessageBodyMD5(resultMessages[i].MessageBody)
			if strings.ToUpper(resultMessages[i].MessageBodyMD5) != want {
				err = fmt.Errorf("the %d'th MessageBodyMD5 mismatch, have:%s, want:%s", i, resultMessages[i].MessageBodyMD5, want)
				return
			}
		}
		if q.config.Base64Enabled {
			resultMessages := result.Messages
			for i := range resultMessages {
				if len(resultMessages[i].MessageBody) == 0 {
					continue
				}
				resultMessages[i].MessageBody, err = internal.Base64Decode(resultMessages[i].MessageBody)
				if err != nil {
					return
				}
				resultMessages[i].MessageBodyMD5 = internal.MessageBodyMD5(resultMessages[i].MessageBody)
			}
		}
		msgs = result.Messages
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

type PeekMessageResponse struct {
	XMLName struct{} `xml:"Message"`

	MessageId        string `xml:"MessageId"`
	MessageBody      []byte `xml:"MessageBody"`
	MessageBodyMD5   string `xml:"MessageBodyMD5"`
	EnqueueTime      int64  `xml:"EnqueueTime"`
	FirstDequeueTime int64  `xml:"FirstDequeueTime"`
	DequeueCount     int    `xml:"DequeueCount"`
	Priority         int    `xml:"Priority"`
}

func (q *Queue) PeekMessage() (requestId string, msg *PeekMessageResponse, err error) {
	return q.PeekMessageContext(context.Background())
}

func (q *Queue) PeekMessageContext(ctx context.Context) (requestId string, msg *PeekMessageResponse, err error) {
	_url, err := internal.ParseURL(q.queue + "/messages?peekonly=true")
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodGet, _url, nil, nil, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result PeekMessageResponse
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		if want := internal.MessageBodyMD5(result.MessageBody); strings.ToUpper(result.MessageBodyMD5) != want {
			err = fmt.Errorf("the MessageBodyMD5 mismatch, have:%s, want:%s", result.MessageBodyMD5, want)
			return
		}
		if q.config.Base64Enabled && len(result.MessageBody) > 0 {
			result.MessageBody, err = internal.Base64Decode(result.MessageBody)
			if err != nil {
				return
			}
			result.MessageBodyMD5 = internal.MessageBodyMD5(result.MessageBody)
		}
		msg = &result
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

func (q *Queue) BatchPeekMessage(numOfMessages int) (requestId string, msgs []PeekMessageResponse, err error) {
	return q.BatchPeekMessageContext(context.Background(), numOfMessages)
}

func (q *Queue) BatchPeekMessageContext(ctx context.Context, numOfMessages int) (requestId string, msgs []PeekMessageResponse, err error) {
	if numOfMessages < 1 || numOfMessages > 16 {
		numOfMessages = 16
	}

	_url, err := internal.ParseURL(q.queue + "/messages?peekonly=true&numOfMessages=" + strconv.Itoa(numOfMessages))
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodGet, _url, nil, nil, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result struct {
			XMLName struct{} `xml:"Messages"`

			Messages []PeekMessageResponse `xml:"Message"`
		}
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		resultMessages := result.Messages
		for i := range resultMessages {
			want := internal.MessageBodyMD5(resultMessages[i].MessageBody)
			if strings.ToUpper(resultMessages[i].MessageBodyMD5) != want {
				err = fmt.Errorf("the %d'th MessageBodyMD5 mismatch, have:%s, want:%s", i, resultMessages[i].MessageBodyMD5, want)
				return
			}
		}
		if q.config.Base64Enabled {
			resultMessages := result.Messages
			for i := range resultMessages {
				if len(resultMessages[i].MessageBody) == 0 {
					continue
				}
				resultMessages[i].MessageBody, err = internal.Base64Decode(resultMessages[i].MessageBody)
				if err != nil {
					return
				}
				resultMessages[i].MessageBodyMD5 = internal.MessageBodyMD5(resultMessages[i].MessageBody)
			}
		}
		msgs = result.Messages
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

func (q *Queue) DeleteMessage(receiptHandle string) (requestId string, err error) {
	return q.DeleteMessageContext(context.Background(), receiptHandle)
}

func (q *Queue) DeleteMessageContext(ctx context.Context, receiptHandle string) (requestId string, err error) {
	_url, err := url.Parse(q.queue + "/messages?ReceiptHandle=" + url.QueryEscape(receiptHandle))
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodDelete, _url, nil, nil, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

type BatchDeleteMessageErrorItem struct {
	XMLName struct{} `xml:"Error"`

	ErrorCode     string `xml:"ErrorCode"`
	ErrorMessage  string `xml:"ErrorMessage"`
	ReceiptHandle string `xml:"ReceiptHandle"`
}

func (q *Queue) BatchDeleteMessage(receiptHandles []string) (requestId string, _errors []BatchDeleteMessageErrorItem, err error) {
	return q.BatchDeleteMessageContext(context.Background(), receiptHandles)
}

func (q *Queue) BatchDeleteMessageContext(ctx context.Context, receiptHandles []string) (requestId string, _errors []BatchDeleteMessageErrorItem, err error) {
	if len(receiptHandles) < 1 || len(receiptHandles) > 16 {
		err = errors.New("the length of receiptHandles is invalid")
		return
	}

	_url, err := internal.ParseURL(q.queue + "/messages")
	if err != nil {
		return
	}

	var req = struct {
		XMLName struct{} `xml:"ReceiptHandles"`

		ReceiptHandles []string `xml:"ReceiptHandle,omitempty"`
	}{
		ReceiptHandles: receiptHandles,
	}

	pool := mns.GetBytesBufferPool()
	reqBuffer := pool.Get()
	defer pool.Put(reqBuffer)
	reqBuffer.Reset()
	if err = xml.NewEncoder(reqBuffer).Encode(&req); err != nil {
		return
	}
	reqBody := reqBuffer.Bytes()

	pool = mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodDelete, _url, nil, reqBody, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		return
	case statusCode == 404:
		if bytes.Contains(respBody, []byte(`</Errors>`)) && bytes.Contains(respBody, []byte(`<ReceiptHandle>`)) { // 部分消息删除失败
			var result struct {
				XMLName struct{} `xml:"Errors"`

				Errors []BatchDeleteMessageErrorItem `xml:"Error"`
			}
			if err = xml.Unmarshal(respBody, &result); err != nil {
				err = internal.NewXMLUnmarshalError(respBody, &result, err)
				return
			}
			_errors = result.Errors
			return
		}
		fallthrough // QueueNotExist
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}

type ChangeMessageVisibilityResponse struct {
	XMLName struct{} `xml:"ChangeVisibility"`

	ReceiptHandle   string `xml:"ReceiptHandle"`
	NextVisibleTime int64  `xml:"NextVisibleTime"`
}

func (q *Queue) ChangeMessageVisibility(receiptHandle string, visibilityTimeout int) (requestId string, resp *ChangeMessageVisibilityResponse, err error) {
	return q.ChangeMessageVisibilityContext(context.Background(), receiptHandle, visibilityTimeout)
}

func (q *Queue) ChangeMessageVisibilityContext(ctx context.Context, receiptHandle string, visibilityTimeout int) (requestId string, resp *ChangeMessageVisibilityResponse, err error) {
	rawurl := q.queue + "/messages?receiptHandle=" + url.QueryEscape(receiptHandle) + "&visibilityTimeout=" + strconv.Itoa(visibilityTimeout)
	_url, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	pool := mns.GetBytesBufferPool()
	respBuffer := pool.Get()
	defer pool.Put(respBuffer)
	respBuffer.Reset()
	requestId, statusCode, respBody, err := internal.DoHTTP(ctx, http.MethodPut, _url, nil, nil, respBuffer, q.config)
	if err != nil {
		return
	}

	switch {
	case statusCode/100 == 2:
		var result ChangeMessageVisibilityResponse
		if err = xml.Unmarshal(respBody, &result); err != nil {
			err = internal.NewXMLUnmarshalError(respBody, &result, err)
			return
		}
		resp = &result
		return
	default:
		err = internal.UnmarshalError(requestId, statusCode, respBody)
		return
	}
}
