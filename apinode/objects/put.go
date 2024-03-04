package objects

import (
	"DisHub/common/utils"
	service "DisHub/service/file_meta"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func put(c *gin.Context) {
	hash := utils.GetHashFromHeader(c)
	if hash == "" {
		c.JSON(http.StatusBadRequest, "miss object hash")
		return
	}
	size := utils.GetSizeFromHeader(c)
	code, err := storeObject(c.Request.Body, hash, size)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to copy request body to file")
		return
	}
	if code != http.StatusOK {
		c.JSON(code, "")
		return
	}
	name := c.Param("name")
	err = service.G_fileMeta.PutMetaData(name, size, hash)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusOK, "")
}
