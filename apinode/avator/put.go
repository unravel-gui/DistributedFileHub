package avator

import (
	"DisHub/apinode/locate"
	"DisHub/common/response"
	"DisHub/common/utils"
	"DisHub/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func put(c *gin.Context) {
	hash := utils.GetHashFromHeader(c)
	if hash == "" {
		response.BadRequest(c, "miss object hash")
		return
	}
	size := utils.GetSizeFromHeader(c)
	// 存在,不重复上传
	if locate.Exist(hash) {
		response.SuccessWithData(c, hash)
		return
	}
	code, err := storeObject(c.Request.Body, hash, size)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to copy request body to file")
		return
	}
	if code != http.StatusOK {
		response.Response(c, code, "upload file failed", nil)
		return
	}
	err = service.G_OssMeta.PutMetaData(hash, size)
	if err != nil {
		log.Println(err)
		response.InternalServer(c, "upload file meta failed")
		return
	}
	response.SuccessWithData(c, hash)
}
