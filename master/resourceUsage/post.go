package resourceUsage

import (
	"DisHub/common/response"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

type PostForm struct {
	Addr string `json:"addr"`
}

func post(c *gin.Context) {
	var form PostForm
	if err := c.BindJSON(&form); err != nil {
		response.BadRequest(c, "parse form err")
		return
	}
	if form.Addr == "" {
		response.BadRequest(c, "addr not be empty")
		return
	}

	fs, err := service.G_ResourceUsage.GetNodeResourceUsage(form.Addr)
	if err != nil {
		response.InternalServer(c, "get resource usage failed")
		return
	}
	response.SuccessWithData(c, fs)
}
