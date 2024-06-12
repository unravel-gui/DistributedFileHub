package user

import (
	"DisHub/common/response"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

func getUserInfo(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	userInfo, err := service.G_User.GetUserInfo(user.Uid)
	if err != nil {
		response.InternalServer(c, "get userInfo failed")
		return
	}
	userInfo.Password = ""
	response.SuccessWithData(c, userInfo)
}
