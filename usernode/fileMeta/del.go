package fileMeta

import (
	"DisHub/common/response"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

type DelFiles struct {
	Fids []int `json:"fids"`
}

func removeFileMetasByUidAndDirs(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.Unauthorized(c, "parse userInfo failed")
		return
	}
	var delFiles DelFiles
	if err := c.BindJSON(&delFiles); err != nil {
		response.InternalServer(c, "parse dirs failed")
		return
	}
	err := service.G_FileMeta.RemoveFileMetasByDir(user.Uid, delFiles.Fids)
	if err != nil {
		response.InternalServer(c, "del file meta failed")
		return
	}
	response.Success(c)
}

func delFileMetasByUidAndHash(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.Unauthorized(c, "parse userInfo failed")
		return
	}
	hash := c.Param("hash")
	if hash == "" {
		response.BadRequest(c, "hash should not be empty")
		return
	}
	ok, err := service.G_FileMeta.DeleteFileMetaByHash(user.Uid, hash)
	if err != nil {
		response.InternalServer(c, "get file meta failed")
		return
	}
	message := "success"
	if !ok {
		message = "file not existed"
	}
	response.SuccessWithMsg(c, message)
}
