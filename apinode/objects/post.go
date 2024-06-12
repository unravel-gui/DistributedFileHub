package objects

import (
	"DisHub/apinode/heartbeat"
	"DisHub/apinode/locate"
	"DisHub/common/response"
	"DisHub/common/rs"
	"DisHub/common/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func post(c *gin.Context) {
	size := utils.GetInt64FromHeader(c, "size")
	hash := utils.GetHashFromHeader(c)
	if hash == "" {
		log.Println("missing object hash in digest header")
		response.BadRequest(c, "missing object hash")
		return
	}
	if locate.Exist(hash) {
		response.SuccessWithMsg(c, "file is existed")
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		response.ServiceUnavailable(c, "cannot find enough dataServer")
		return
	}
	stream, e := rs.NewRSResumablePutStream(ds, hash, size)
	if e != nil {
		log.Println(e)
		response.InternalServerByError(c, e)
		return
	}
	token := stream.ToToken()
	c.Writer.Header().Set("token", token)
	response.SuccessWithData(c, struct {
		Token string `json:"token"`
	}{Token: token})
}
