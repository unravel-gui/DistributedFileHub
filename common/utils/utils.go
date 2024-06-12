package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
	"time"
)

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func CalculateStringHash(str string) string {
	r := strings.NewReader(str)
	return CalculateHash(r)
}

// 获取当前时间，并设置时区为北京时间
var loc, _ = time.LoadLocation("Asia/Shanghai")

func GetNow() time.Time {
	now := time.Now().In(loc)
	return now
}
