package objects

import (
	"DisHub/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func put(c *gin.Context) {
	filename := c.Param("name")
	filePath := config.GetFilePath(filename)
	f, e := os.Create(filePath)
	if e != nil {
		log.Println(e)
		c.JSON(http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer f.Close()
	_, err := io.Copy(f, c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to copy request body to file")
		return
	}
	log.Println("ok")
	c.JSON(http.StatusOK, "success")
}
