package fileMeta

import (
	"DisHub/middleware"
	"DisHub/repository"
	"github.com/gin-gonic/gin"
)

func Handler(path string, r *gin.Engine) {
	router := r.Group(path)
	router.Use(middleware.CheckJWTToken())
	// 通过uid和dir获得文件是元数据
	router.GET("/baseFolder", get)
	router.GET("/folder/:dir", getFileMetasByUidAndDir)
	router.GET("/delFolder/:dir", getDelFileMetasByUidAndDir)
	router.GET("/hash/:hash", getFileMetasByUidAndHash)
	router.POST("/folder", postFileMetas)
	router.POST("/folder/recover", recoverFileMetasByUidAndDirs)
	router.POST("/folder/del", delFileMetasByUidAndDirs)
	router.PUT("/folder", putFileMeta)
	router.DELETE("/folder/remove", removeFileMetasByUidAndDirs)
	router.DELETE("/hash/:hash", delFileMetasByUidAndHash)
}

func GetUserInfoFromContext(c *gin.Context) *repository.User {
	userS, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := userS.(*repository.User)
	if !ok {
		return nil
	}
	return user
}
