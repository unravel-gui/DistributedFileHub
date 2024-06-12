package user

import (
	"DisHub/common/response"
	"DisHub/repository"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

func putUserInfo(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	var userReq repository.User
	// 解析请求体内容到 user 结构体中
	if err := c.BindJSON(&userReq); err != nil {
		response.BadRequestByError(c, err)
		return
	}
	user, err := service.G_User.UpdateUserInfo(user.Uid, &userReq)
	if err != nil {
		response.InternalServer(c, "create user failed")
		return
	}
	user.Password = ""
	response.SuccessWithData(c, user)
}

type ChangePass struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

func putUserPassword(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	var userReq ChangePass
	// 解析请求体内容到 user 结构体中
	if err := c.BindJSON(&userReq); err != nil {
		response.BadRequestByError(c, err)
		return
	}
	ok, err := service.G_User.UpdateUserPassword(user.Uid, userReq.OldPass, userReq.NewPass)
	if !ok {
		response.BadRequest(c, "old password is wrong")
		return
	}
	if err != nil {
		response.InternalServer(c, "update password failed")
		return
	}
	response.Success(c)
}
