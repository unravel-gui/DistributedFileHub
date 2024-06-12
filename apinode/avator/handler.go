package avator

import "github.com/gin-gonic/gin"

func Handler(path string, c *gin.Engine) {
	router := c.Group(path)
	router.GET("/:hash", get)
	router.PUT("/avator", put)
}
