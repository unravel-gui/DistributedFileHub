package temp

import (
	"DisHub/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func patch(c *gin.Context) {
	uuid := c.Param("uuid")
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusNotFound, "")
		return
	}
	infoFile := config.GetBasePath() + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	defer f.Close()
	_, e = io.Copy(f, c.Request.Body)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	actual := info.Size()
	//
	log.Println("current size = ", actual)
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		c.JSON(http.StatusInternalServerError, "")
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	f, e := os.Open(config.GetBasePath() + "/temp/" + uuid)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}
