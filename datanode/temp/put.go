package temp

import (
	"DisHub/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func put(c *gin.Context) {
	uuid := c.Param("uuid")
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusNotFound, "not found")
		return
	}
	infoFile := config.GetBasePath() + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		f.Close()
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	e = f.Close()
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		c.JSON(http.StatusInternalServerError, "")
		f.Close()
		return
	}
	commitTempObject(datFile, tempinfo)
	c.JSON(http.StatusOK, "")
}
