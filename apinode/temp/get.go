package temp

import (
	"DisHub/common/response"
	"DisHub/common/rs"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func get(c *gin.Context) {
	token := c.Param("token")
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		response.Forbidden(c, "invalided token")
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		response.Forbidden(c, "get file size failed")
		return
	}
	c.Writer.Header().Set("current_size", fmt.Sprintf("%d", current))
	response.SuccessWithData(c, struct {
		CurrentSize int64 `json:"current_size"`
	}{CurrentSize: current})
}
