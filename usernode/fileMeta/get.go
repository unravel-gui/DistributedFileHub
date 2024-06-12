package fileMeta

import (
	"DisHub/common"
	"DisHub/common/response"
	"DisHub/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FidsResp response.FidsResp

// 获得用户的根目录文件列表
func get(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	fms, err := service.G_FileMeta.GetRootFolder(user.Uid)
	if err != nil || len(fms) < 3 {
		response.InternalServer(c, "get file meta failed")
		return
	}
	rresp := FidsResp{
		HomeFolder:  response.NewFolderInfo(&fms[0]),
		VideoFolder: response.NewFolderInfo(&fms[1]),
		ImageFolder: response.NewFolderInfo(&fms[2]),
	}
	response.SuccessWithData(c, rresp)
}

func getFileMetasByUidAndDir(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	dirStr := c.Param("dir")
	dir, err := strconv.Atoi(dirStr)
	if err != nil {
		response.BadRequest(c, "parse dir failed")
		return
	}
	fms, err := service.G_FileMeta.GetFileMetasByUserAndDir(user.Uid, dir)
	if err != nil {
		response.InternalServer(c, "get file meta failed")
		return
	}
	response.SuccessWithData(c, fms)
}

func getDelFileMetasByUidAndDir(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	dirStr := c.Param("dir")
	dir, err := strconv.Atoi(dirStr)
	if err != nil {
		response.BadRequest(c, "parse dir failed")
		return
	}
	if dir <= 0 {
		dir = common.ROOTFOLDER
	}
	fms, err := service.G_FileMeta.GetDeleteFileMetasByUserAndDir(user.Uid, dir)
	if err != nil {
		response.InternalServer(c, "get file meta failed")
		return
	}
	response.SuccessWithData(c, fms)
}

func getFileMetasByUidAndHash(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	hash := c.Param("hash")
	if hash == "" {
		response.BadRequest(c, "hash should not be empty")
		return
	}
	fm, err := service.G_FileMeta.GetFileMetaByHash(user.Uid, hash)
	if err != nil {
		response.InternalServer(c, "get file meta failed")
		return
	}
	response.SuccessWithData(c, fm)
}
