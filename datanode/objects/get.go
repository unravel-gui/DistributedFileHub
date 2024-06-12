package objects

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func get(c *gin.Context) {
	filename := c.Param("name")
	file := getFile(filename)
	if file == "" {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("open file err:%v", err))
		return
	}
	defer f.Close()
	if _, err = io.Copy(c.Writer, f); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to copy file content to response")
		return
	}
}
