package heartbeat

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/config"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()

	for {
		q.Publish(common.EXCHANGE_API, config.GetLocalAddr())
		time.Sleep(5 * time.Second)
	}
}
