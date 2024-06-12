package resourceUsage

import (
	"DisHub/common/response"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	rss := service.G_ResourceUsage.GetResourceLastest()
	response.SuccessWithData(c, rss)
}
