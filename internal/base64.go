package internal

import (
	"encoding/base64"
	"fmt"
)

func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func Base64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, src=%s, error=%s", src, err.Error())
	}
	return dst[:n], nil
}
