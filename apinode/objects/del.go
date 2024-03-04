package objects

import (
	service "DisHub/service/file_meta"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func del(c *gin.Context) {
	hash := c.Param("name") // hash
	err := service.G_fileMeta.DelMetaData(hash)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusOK, "")
}
