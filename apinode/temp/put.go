package temp

import (
	"DisHub/apinode/locate"
	"DisHub/common/response"
	"DisHub/common/rs"
	"DisHub/common/utils"
	"DisHub/service"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func put(c *gin.Context) {
	token := c.Param("token")
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		response.InternalServer(c, "")
		return
	}
	// 获得当前文件大小
	current := stream.CurrentSize()
	uploadedSize := current
	if current == -1 {
		response.Forbidden(c, "get file size failed")
		return
	}
	// 获得偏移值
	offset := utils.GetOffsetFromHeader(c)
	if current != offset {
		response.Response(c, http.StatusPartialContent, "offset not match", struct {
			UploadedSize int64 `json:"uploaded_size"`
		}{UploadedSize: uploadedSize})
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, e := io.ReadFull(c.Request.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			response.InternalServer(c, "read body failed")
			return
		}
		current += int64(n)
		if current > stream.Size {
			// 长度异常
			stream.Commit(false)
			log.Println("resumable put exceed size")
			response.Forbidden(c, "resumable put exceed size")
			return
		}
		// 数据检查
		if n != rs.BLOCK_SIZE && current != stream.Size {
			response.Response(c, http.StatusPartialContent, "", struct {
				UploadedSize int64 `json:"uploaded_size"`
			}{UploadedSize: uploadedSize})
			return
		}
		// 当前分块数据上传datanode
		stream.Write(bytes[:n])
		uploadedSize = uploadedSize + int64(n)
		if current == stream.Size {
			// 纠删码处理后上传
			stream.Flush()
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if e != nil {
				stream.Commit(false)
				log.Println("getStream err", e)
				response.Forbidden(c, "commit failed")
				return
			}
			hash := utils.CalculateHash(getStream)
			// 校验Hash
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				response.Forbidden(c, "hash not match")
				return
			}
			if locate.Exist(hash) {
				// 避免上传完成
				stream.Commit(false)
				return
			}
			e = service.G_OssMeta.PutMetaData(stream.Hash, stream.Size)
			if e != nil {
				stream.Commit(false)
				log.Println(e)
				response.InternalServer(c, "upload file meta failed")
				return
			}
			stream.Commit(true)
			response.Response(c, http.StatusPartialContent, "", struct {
				UploadedSize int64 `json:"uploaded_size"`
			}{UploadedSize: uploadedSize})
			return
		}
	}
}
