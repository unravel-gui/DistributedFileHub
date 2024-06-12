package fileMeta

import (
	"DisHub/common"
	"DisHub/common/response"
	"DisHub/repository"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

func putFileMeta(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	var pfm *repository.FileMeta
	if err := c.BindJSON(&pfm); err != nil {
		response.BadRequest(c, "parse file meta err")
		return
	}
	if pfm.Dir == common.ROOTFOLDER {
		response.BadRequest(c, "bad folder")
		return
	}
	ok, err := service.G_FileMeta.CheckUserFileOwnership(user.Uid, pfm.Dir)
	if err != nil {
		response.InternalServerByError(c, err)
		return
	}
	if !ok {
		response.Unauthorized(c, "wrong parent folder")
		return
	}
	ok, err = service.G_FileMeta.CheckUserFileExisted(user.Uid, pfm.Dir, pfm.Name)
	if err != nil {
		response.InternalServerByError(c, err)
		return
	}
	if ok {
		response.BadRequest(c, "File name duplicate")
		return
	}
	ok, err = service.G_FileMeta.UpdateFileMeta(user.Uid, pfm)
	if err != nil || !ok {
		response.InternalServer(c, "update file meta failed")
		return
	}
	response.Success(c)
}
