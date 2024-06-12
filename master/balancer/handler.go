package balancer

import "github.com/gin-gonic/gin"

func Handler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.GET("/:nodeType", get)
}
