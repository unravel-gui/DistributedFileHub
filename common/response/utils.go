package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context) {
	Response(c, http.StatusOK, "success", nil)
}

func SuccessWithMsg(c *gin.Context, msg string) {
	Response(c, http.StatusOK, msg, nil)
}

func SuccessWithData(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, "success", data)
}

func BadRequest(c *gin.Context, msg string) {
	Response(c, http.StatusBadRequest, msg, nil)
}
func BadRequestByError(c *gin.Context, err error) {
	Response(c, http.StatusBadRequest, err.Error(), nil)
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	})
	Response(c, http.StatusUnauthorized, msg, nil)
}

func InternalServer(c *gin.Context, msg string) {
	Response(c, http.StatusInternalServerError, msg, nil)
}

func InternalServerByError(c *gin.Context, err error) {
	Response(c, http.StatusServiceUnavailable, err.Error(), nil)
}

func ServiceUnavailable(c *gin.Context, msg string) {
	Response(c, http.StatusServiceUnavailable, msg, nil)
}

func Forbidden(c *gin.Context, msg string) {
	Response(c, http.StatusForbidden, msg, nil)
}

func NotFound(c *gin.Context, msg string) {
	Response(c, http.StatusNotFound, msg, nil)
}

func Response(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
