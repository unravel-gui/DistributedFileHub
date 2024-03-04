package objects

import (
	"DisHub/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func get(c *gin.Context) {
	filename := c.Param("name")
	filePath := config.GetFilePath(filename)
	f, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		c.String(http.StatusNotFound, "File not found")
		return
	}
	defer f.Close()
	if _, err = io.Copy(c.Writer, f); err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Failed to copy file content to response")
		return
	}
}
