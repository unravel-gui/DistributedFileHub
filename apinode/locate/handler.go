package locate

import "github.com/gin-gonic/gin"

func Handler(path string, c *gin.Engine) {
	router := c.Group(path)
	router.GET("/:name", locate)
}
