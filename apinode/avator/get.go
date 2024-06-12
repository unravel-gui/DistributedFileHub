package avator

import (
	"DisHub/apinode/objects"
	"DisHub/common/response"
	"DisHub/common/utils"
	"DisHub/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func get(c *gin.Context) {
	hash := c.Param("hash") // hash
	meta, e := service.G_OssMeta.GetMetaData(hash)
	if e != nil {
		log.Println(e)
		response.InternalServer(c, "get file meta err")
		return
	}
	if meta == nil {
		response.NotFound(c, "File not found")
		return
	}

	stream, err := objects.GetStream(hash, meta.Size)
	if err != nil {
		log.Println(err)
		response.NotFound(c, "File not found")
		return
	}
	defer stream.Close()
	// 移到指定offset
	offset := utils.GetOffsetFromHeader(c)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		c.Writer.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		c.JSON(http.StatusPartialContent, "")
	}
	// 设置请求头
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", meta.Hash))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	if _, err = io.Copy(c.Writer, stream); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
