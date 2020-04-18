package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

/* 后台登录相关 */
const (
	webkey = "UXnezMezvpLnQre6XvhdWpHxiWFSrBZR"
)

func getPassword(src string) (string, error) {
	src += webkey

	h := sha256.New()
	h.Write([]byte(src))
	sum := h.Sum(nil)
	s := hex.EncodeToString(sum)

	h1 := md5.New()
	h1.Write([]byte(s))
	s = hex.EncodeToString(h1.Sum(nil))

	return s, nil
}
