package temp

import (
	"DisHub/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func head(c *gin.Context) {
	uuid := c.Param("uuid")
	f, e := os.Open(config.GetBasePath() + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusNotFound, "")
		return
	}
	defer f.Close()
	// 文件信息
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	// 文件大小
	c.Writer.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
