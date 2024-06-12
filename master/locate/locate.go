package locate

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/config"
	"fmt"
)

func StartLocate() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()
	q.Bind(common.EXCHANGE_LOCATE_MASTER)
	c := q.Consume()
	for msg := range c {
		fmt.Println(msg.ReplyTo)
		q.Send(msg.ReplyTo, config.GetLocalAddr())
	}
}
