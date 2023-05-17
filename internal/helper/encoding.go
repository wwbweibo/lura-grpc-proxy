package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func Base64(bts []byte) string {
	return base64.StdEncoding.EncodeToString(bts)
}

func HmacHd5(str string, key string) string {
	sig := hmac.New(md5.New, []byte(key))
	sig.Write([]byte(str))
	return hex.EncodeToString(sig.Sum(nil))
}
