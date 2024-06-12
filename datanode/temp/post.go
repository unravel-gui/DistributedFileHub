package temp

import (
	"DisHub/common/utils"
	"DisHub/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func post(c *gin.Context) {
	uuidStr := uuid.New().String()
	name := c.Param("uuid")
	size := utils.GetInt64FromHeader(c, "size")
	//
	log.Println("temp file Size = ", size)
	t := tempInfo{uuidStr, name, size}
	e := t.writeToFile()
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	f, _ := os.Create(config.GetBasePath() + "/temp/" + t.Uuid + ".dat")
	f.Close()
	c.String(http.StatusOK, uuidStr)
}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(config.GetBasePath() + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}
