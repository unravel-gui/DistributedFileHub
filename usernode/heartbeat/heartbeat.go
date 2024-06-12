package heartbeat

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/common/resource"
	"DisHub/common/types"
	"DisHub/config"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()

	for {
		q.Publish(common.EXCHANGE_MASTER, types.HeartBeatMessage{
			Addr:           config.GetLocalAddr(),
			Weight:         config.GetNodeWeight(),
			NodeType:       common.USERNODE,
			ResourceStatus: resource.NewResourceStatus(),
		})
		time.Sleep(common.HEATBEAT_INTERVAL * time.Second)
	}
}
