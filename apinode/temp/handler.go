package temp

import (
	"DisHub/repository"
	"github.com/gin-gonic/gin"
)

func Handler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.PUT("/:token", put)
	router.GET("/:token", get)
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
