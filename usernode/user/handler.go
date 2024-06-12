package user

import (
	"DisHub/repository"
	"github.com/gin-gonic/gin"
)

func Handler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.GET("/info", getUserInfo)
	router.PUT("/info", putUserInfo)
	router.PUT("/password", putUserPassword)
}

func HandlerWithoutCheck(path string, r *gin.Engine) {
	router := r.Group(path)
	router.POST("/register", register)
	router.POST("/login", login)
}

func GetUserInfoFromContext(c *gin.Context) *repository.User {
	userS, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := userS.(*repository.User)
	if !ok {
		return nil
	}
	return user
}
