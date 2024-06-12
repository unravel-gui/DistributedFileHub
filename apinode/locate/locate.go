package locate

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/common/response"
	"DisHub/common/rs"
	"DisHub/common/types"
	"DisHub/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func get(c *gin.Context) {
	hash := c.Param("hash")
	info := Locate(hash)
	if len(info) == 0 {
		response.SuccessWithData(c, struct {
			Code int `json:"code"`
		}{
			Code: http.StatusNotFound,
		})
		return
	}
	response.Success(c)
}

func locate(c *gin.Context) {
	hash := c.Param("hash")
	info := Locate(hash)
	if len(info) == 0 {
		response.NotFound(c, "file not found")
		return
	}
	c.JSON(http.StatusOK, info)
}

func Locate(hash string) (locateInfo map[int]string) {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	q.Publish(common.EXCHANGE_DATA, hash)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
