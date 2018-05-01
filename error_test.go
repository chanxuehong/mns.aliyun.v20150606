package mns

import "testing"

func TestIsMessageNotExist(t *testing.T) {
	have := IsMessageNotExist(nil)
	want := false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	var err *Error
	have = IsMessageNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{}
	have = IsMessageNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeMessageNotExist,
		Code:           ErrorCodeQueueNotExist,
	}
	have = IsMessageNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeMessageNotExist + 1,
		Code:           ErrorCodeMessageNotExist,
	}
	have = IsMessageNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeMessageNotExist,
		Code:           ErrorCodeMessageNotExist,
	}
	have = IsMessageNotExist(err)
	want = true
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}
}

func TestIsReceiptHandleError(t *testing.T) {
	have := IsReceiptHandleError(nil)
	want := false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	var err *Error
	have = IsReceiptHandleError(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{}
	have = IsReceiptHandleError(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeReceiptHandleError,
		Code:           ErrorCodeQueueNotExist,
	}
	have = IsReceiptHandleError(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeReceiptHandleError + 1,
		Code:           ErrorCodeReceiptHandleError,
	}
	have = IsReceiptHandleError(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeReceiptHandleError,
		Code:           ErrorCodeReceiptHandleError,
	}
	have = IsReceiptHandleError(err)
	want = true
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}
}

func TestIsQueueNotExist(t *testing.T) {
	have := IsQueueNotExist(nil)
	want := false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	var err *Error
	have = IsQueueNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{}
	have = IsQueueNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeQueueNotExist,
		Code:           ErrorCodeTopicNotExist,
	}
	have = IsQueueNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeQueueNotExist + 1,
		Code:           ErrorCodeQueueNotExist,
	}
	have = IsQueueNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeQueueNotExist,
		Code:           ErrorCodeQueueNotExist,
	}
	have = IsQueueNotExist(err)
	want = true
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}
}

func TestIsTopicNotExist(t *testing.T) {
	have := IsTopicNotExist(nil)
	want := false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	var err *Error
	have = IsTopicNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{}
	have = IsTopicNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeTopicNotExist,
		Code:           ErrorCodeQueueNotExist,
	}
	have = IsTopicNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeTopicNotExist + 1,
		Code:           ErrorCodeTopicNotExist,
	}
	have = IsTopicNotExist(err)
	want = false
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}

	err = &Error{
		HttpStatusCode: ErrorHttpStatusCodeTopicNotExist,
		Code:           ErrorCodeTopicNotExist,
	}
	have = IsTopicNotExist(err)
	want = true
	if have != want {
		t.Errorf("have:%t, want:%t", have, want)
		return
	}
}
