package temp

import (
	"github.com/gin-gonic/gin"
)

func Handler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.POST("/:hash", post)
	router.PUT("/:uuid", put)
	router.PATCH("/:uuid", patch)
	router.DELETE("/:uuid", del)
}
