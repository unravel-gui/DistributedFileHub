package locate

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// var rabbitmqServer =os.Getenv("RABBITMQ_SERVER")
var rabbitmqServer = "amqp://osstest:osstest@192.168.52.140:5672"

func locate(c *gin.Context) {
	name := c.Param("name")
	info := Locate(name)
	if len(info) == 0 {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	c.JSON(http.StatusOK, info)
}

func Locate(name string) string {
	q := rabbitmq.New(rabbitmqServer)
	q.Publish(common.EXCHANGE_DATA, name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
