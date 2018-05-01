package internal

import (
	"net/url"
	"testing"
)

func TestParseURL(t *testing.T) {
	u1, err := ParseURL("https://www.google.com/test?ka=va&ka=va2&kb=vb")
	if err != nil {
		t.Error(err.Error())
		return
	}
	u2, err := url.Parse("https://www.google.com/test?ka=va&ka=va2&kb=vb")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if *u1 != *u2 || u1.RequestURI() != u2.RequestURI() {
		t.Errorf("have:%v, want:%v", u1, u2)
		return
	}
}

func BenchmarkParseURL(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ParseURL("https://www.google.com/test?ka=va&ka=va2&kb=vb"); err != nil {
			b.Error(err.Error())
			return
		}
	}
}

func BenchmarkURLParse(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := url.Parse("https://www.google.com/test?ka=va&ka=va2&kb=vb"); err != nil {
			b.Error(err.Error())
			return
		}
	}
}
