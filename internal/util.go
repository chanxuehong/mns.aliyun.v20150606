package internal

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"time"
)

// Authorization returns Authorization http header.
func Authorization(accessKeyId, signature string) string {
	return "MNS " + accessKeyId + ":" + signature
}

// "Thu, 17 Mar 2012 18:49:58 GMT"
func FormatDate(t time.Time) string {
	return t.UTC().Format(DateFormatLayout)
}

// ContentMD5 returns md5 of body.
//  https://tools.ietf.org/html/rfc1864
func ContentMD5(b []byte) string {
	sum := md5.Sum(b)
	return base64.StdEncoding.EncodeToString(sum[:])
}

// MessageBodyMD5 returns md5 of MessageBody.
func MessageBodyMD5(b []byte) string {
	sum := md5.Sum(b)
	var hexSum [md5.Size * 2]byte
	hex.Encode(hexSum[:], sum[:])
	return string(bytes.ToUpper(hexSum[:]))
}
