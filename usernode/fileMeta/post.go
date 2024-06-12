package fileMeta

import (
	"DisHub/common"
	"DisHub/common/response"
	"DisHub/common/utils"
	"DisHub/repository"
	"DisHub/service"
	"github.com/gin-gonic/gin"
)

type PostFileMeta struct {
	Dir         *int   `json:"dir"`
	Hash        string `json:"hash"`
	Name        string `json:"name"`
	Size        *int64 `json:"size"`
	ContentType string `json:"content_type"`
	Force       bool   `json:"force"`
}

func postFileMetas(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	var pfm PostFileMeta
	if err := c.BindJSON(&pfm); err != nil {
		response.BadRequest(c, "parse file meta err")
		return
	}
	ok, msg := pfm.Check()
	if !ok {
		response.BadRequest(c, msg)
		return
	}

	ok, err := service.G_FileMeta.CheckUserFileOwnership(user.Uid, *pfm.Dir)
	if err != nil {
		response.InternalServerByError(c, err)
		return
	}
	if !ok {
		response.Unauthorized(c, "wrong parent folder")
		return
	}
	fm := newFileMeta(&pfm, user)
	// 非强制则提前检查是否存在重复
	if !pfm.Force {
		ok, err = service.G_FileMeta.CheckUserFileExisted(user.Uid, *pfm.Dir, pfm.Name)
		if err != nil {
			response.InternalServerByError(c, err)
			return
		}
		if ok {
			response.BadRequest(c, "File name duplicate")
			return
		}
		ok, err = service.G_FileMeta.PutFileMeta(fm)
		if err != nil {
			response.InternalServer(c, "put file meta failed")
			return
		}
		if !ok {
			response.SuccessWithMsg(c, "file is existed")
			return
		}
	} else {
		ok, err = service.G_FileMeta.PutFileMetaForce(fm)
		if err != nil {
			response.InternalServer(c, "put file meta failed")
			return
		}
		if !ok {
			response.SuccessWithMsg(c, "file is existed")
			return
		}
	}
	response.Success(c)
}

type ParamFiles struct {
	Fids []int `json:"fids"`
}

func recoverFileMetasByUidAndDirs(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.BadRequest(c, "parse userInfo failed")
		return
	}
	var recoverFiles ParamFiles
	if err := c.BindJSON(&recoverFiles); err != nil {
		response.BadRequest(c, "parse dirs failed")
		return
	}
	//TODO:recover
	err := service.G_FileMeta.LogicalRecoverFileMetasByDir(user.Uid, recoverFiles.Fids)
	if err != nil {
		response.InternalServer(c, "recover file meta failed")
		return
	}
	response.Success(c)
}

func delFileMetasByUidAndDirs(c *gin.Context) {
	user := GetUserInfoFromContext(c)
	if user == nil {
		response.Unauthorized(c, "parse userInfo failed")
		return
	}
	var delFiles ParamFiles
	if err := c.BindJSON(&delFiles); err != nil {
		response.InternalServer(c, "parse dirs failed")
		return
	}
	err := service.G_FileMeta.LogicalDeleteFileMetasByDir(user.Uid, delFiles.Fids)
	if err != nil {
		response.InternalServer(c, "del file meta failed")
		return
	}
	response.Success(c)
}
func (pfm *PostFileMeta) Check() (bool, string) {
	msg := ""
	if pfm.Dir == nil {
		msg += "dir "
	}
	if pfm.Size == nil {
		msg += "size "
	}
	if pfm.Name == "" {
		msg += "name "
	}
	if pfm.Hash == "" {
		msg += "hash "
	}
	if pfm.ContentType == "" {
		msg += "content_type "
	}

	if len(msg) != 0 {
		msg += "should not be empty"
		return false, msg
	}
	if *pfm.Dir == common.ROOTFOLDER {
		msg = "parent folder not less 0"
		return false, msg
	}
	return true, ""
}

func newFileMeta(pfm *PostFileMeta, user *repository.User) *repository.FileMeta {
	now := utils.GetNow()
	fm := &repository.FileMeta{
		Uid:         user.Uid,
		Name:        pfm.Name,
		Hash:        pfm.Hash,
		Dir:         *pfm.Dir,
		Size:        *pfm.Size,
		ContentType: pfm.ContentType,
		UploadTime:  now,
		UpdateTime:  now,
	}
	return fm
}
