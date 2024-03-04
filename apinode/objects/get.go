package objects

import (
	"DisHub/common/utils"
	service "DisHub/service/file_meta"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func get(c *gin.Context) {
	name := c.Param("name") // hash
	hash := utils.GetHashFromHeader(c)
	meta, e := service.G_fileMeta.GetMetaData(hash)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "get file meta err")
		return
	}
	if meta == nil {
		c.JSON(http.StatusNotFound, "not found")
		return
	}

	stream, err := getStream(name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, "File not found")
		return
	}
	if _, err = io.Copy(c.Writer, stream); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to copy file content to response")
		return
	}
}
