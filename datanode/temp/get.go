package temp

import (
	"DisHub/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func get(c *gin.Context) {
	uuid := c.Param("uuid")
	f, e := os.Open(config.GetBasePath() + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusNotFound, "")
		return
	}
	defer f.Close()
	io.Copy(c.Writer, f)
}
