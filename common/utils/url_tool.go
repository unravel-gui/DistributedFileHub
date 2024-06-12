package utils

import (
	"DisHub/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetHashFromHeader(c *gin.Context) string {
	return c.Request.Header.Get("Hash")
}

func SetMagicTokenFromHeader(h http.Header) {
	h.Set(common.RPC_TOKEN_KEY, common.RPC_TOKEN_VALVE)
}

func SetJWTTokenFromHeader(c *gin.Context, token string) {
	c.Header(common.JWT_TOKEN, token)
}

func GetSizeFromHttpHeader(h http.Header) int64 {
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	return size
}

func GetSizeFromHeader(c *gin.Context) int64 {
	return GetInt64FromHeader(c, "Content-Length")
}
func GetInt64FromHeader(c *gin.Context, key string) int64 {
	length := c.Request.Header.Get(key)
	size, _ := strconv.ParseInt(length, 0, 64)
	return size
}
func GetOffsetFromHeader(c *gin.Context) int64 {
	byteRange := c.Request.Header.Get("Content-Range")
	// 验证参数
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}
