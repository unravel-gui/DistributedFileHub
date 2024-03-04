package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetHashFromHeader(c *gin.Context) string {
	return c.Request.Header.Get("Hash")
}

func GetSizeFromHeader(c *gin.Context) int64 {
	return GetInt64FromHeader(c, "Content-Length")
}

func GetInt64FromHeader(c *gin.Context, key string) int64 {
	length := c.Request.Header.Get(key)
	size, _ := strconv.ParseInt(length, 0, 64)
	return size
}
