package proxy

import (
	"github.com/gin-gonic/gin"
)

func UserHandler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.Any("/*path", userProxy)
}
func UserHandlerWithOutCheck(path string, r *gin.Engine) {
	router := r.Group(path)
	router.POST("/register", userProxy)
	router.POST("/login", userProxy)
}

func MasterHandler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.Any("/*path", masterProxy)
}
