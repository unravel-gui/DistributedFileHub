package user

import (
	"DisHub/common"
	"DisHub/common/jwt"
	"DisHub/common/response"
	"DisHub/common/utils"
	"DisHub/repository"
	"DisHub/service"
	"github.com/gin-gonic/gin"
	"time"
)

type FidsResp response.FidsResp
type RegisterUser struct {
	repository.User
	Code string `json:"code"`
}

func register(c *gin.Context) {
	var user RegisterUser
	// 解析请求体内容到 user 结构体中
	if err := c.BindJSON(&user); err != nil {
		response.BadRequestByError(c, err)
		return
	}
	user.Nickname = user.Username
	user.User.Avatar = "https://pic.imeitou.com/uploads/allimg/220908/8-220ZQH328-50.jpg"
	if user.Code == "admin" {
		user.User.IsAdmin = 1
	} else {
		user.User.IsAdmin = 0
	}
	ok, err := service.G_User.Register(&user.User)
	if err != nil {
		response.InternalServer(c, "create user failed")
		return
	}
	if !ok {
		response.BadRequest(c, "user is existed")
		return
	}
	// 初始化文件夹
	folders := initFolder(&user.User)
	ok, err = service.G_FileMeta.BatchPutFileMeta(folders)
	if !ok {
		service.G_User.Remove(&user.User)
		response.InternalServer(c, "init folder failed")
		return
	}
	token, err := jwt.GenerateUserToken(&user.User)
	if err != nil {
		response.InternalServer(c, "generate user token failed")
		return
	}
	utils.SetJWTTokenFromHeader(c, token)
	rresp := FidsResp{
		JwtToken:    token,
		HomeFolder:  response.NewFolderInfo((*folders)[0]),
		VideoFolder: response.NewFolderInfo((*folders)[1]),
		ImageFolder: response.NewFolderInfo((*folders)[2]),
	}
	response.SuccessWithData(c, rresp)
}

func login(c *gin.Context) {
	var user *repository.User
	// 解析请求体内容到 user 结构体中
	if err := c.BindJSON(&user); err != nil {
		response.BadRequestByError(c, err)
		return
	}
	user, err := service.G_User.Login(user)
	if err != nil {
		response.InternalServer(c, "get user login info err")
		return
	}
	if user == nil {
		response.BadRequest(c, "username or password is wrong")
		return
	}
	fms, err := service.G_FileMeta.GetRootFolder(user.Uid)
	if err != nil || len(fms) < 3 {
		response.InternalServer(c, "get file meta failed")
		return
	}
	token, err := jwt.GenerateUserToken(user)
	if err != nil {
		response.InternalServer(c, "generate user token failed")
		return
	}
	utils.SetJWTTokenFromHeader(c, token)
	rresp := FidsResp{
		JwtToken:    token,
		HomeFolder:  response.NewFolderInfo(&fms[0]),
		VideoFolder: response.NewFolderInfo(&fms[1]),
		ImageFolder: response.NewFolderInfo(&fms[2]),
	}
	response.SuccessWithData(c, rresp)
}

func initFolder(user *repository.User) *[]*repository.FileMeta {
	now := utils.GetNow()
	fms := []*repository.FileMeta{
		newFolder(common.REGULAR, user.Uid, now),
		newFolder(common.VIDEOS, user.Uid, now),
		newFolder(common.IMAGES, user.Uid, now),
	}
	return &fms
}
func newFolder(folderName string, uid int, now time.Time) *repository.FileMeta {
	return &repository.FileMeta{
		Uid:         uid,
		Dir:         common.ROOTFOLDER,
		Name:        folderName,
		ContentType: common.FOLDER,
		UploadTime:  now,
		UpdateTime:  now,
	}
}
