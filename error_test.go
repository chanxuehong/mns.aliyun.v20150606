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
