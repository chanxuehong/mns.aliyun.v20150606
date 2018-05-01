package internal

import (
	"testing"
	"time"
)

func TestAuthorization(t *testing.T) {
	have := Authorization("accessKeyId", "signature")
	want := "MNS accessKeyId:signature"
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}

func TestFormatDate(t *testing.T) {
	have := FormatDate(time.Date(2018, 5, 1, 15, 4, 5, 0, time.UTC))
	want := "Tue, 01 May 2018 15:04:05 GMT"
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}

func TestContentMD5(t *testing.T) {
	have := ContentMD5([]byte("1234567890"))
	want := "6Afx/PgtEy+bsBjKZzihnw=="
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}

func TestMessageBodyMD5(t *testing.T) {
	have := MessageBodyMD5([]byte("1234567890"))
	want := "E807F1FCF82D132F9BB018CA6738A19F"
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}
