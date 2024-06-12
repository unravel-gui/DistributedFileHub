package heartbeat

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/common/resource"
	"DisHub/common/types"
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/loadbalancer"
	"DisHub/service"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

var NodeMap = make(map[common.NodeType]*sync.Map)
var lbm = loadbalancer.G_LoadBalancerMap

func registerMap() {
	// 赋值为sync.map
	NodeMap[common.MASTER] = &sync.Map{}
	NodeMap[common.DATANODE] = &sync.Map{}
	NodeMap[common.APINODE] = &sync.Map{}
	NodeMap[common.USERNODE] = &sync.Map{}
}

func ListenHeartbeat() {
	registerMap()
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()
	q.Bind(common.EXCHANGE_MASTER)
	c := q.Consume()
	for k, _ := range NodeMap {
		go removeExpiredServer(k)
	}
	for msg := range c {
		if len(msg.Body) == 0 {
			return
		}
		var info types.HeartBeatMessage
		json.Unmarshal(msg.Body, &info)
		serverMap, existed := NodeMap[info.NodeType]
		fmt.Println(info)
		service.G_ResourceUsage.InsertOne(info)
		if !existed {
			log.Println("NodeType is not existed")
			continue
		}
		_, loaded := serverMap.Load(info.Addr)
		if !loaded {
			lbm.AddNode(info.NodeType, info.Addr, info.Weight)
		}
		// 更新serverMap中的对应id的时间
		serverMap.Store(info.Addr, utils.GetNow())
	}
}

func removeExpiredServer(t common.NodeType) {
	// 取出对应的servers
	serverMap, ok := NodeMap[t]
	if !ok {
		log.Println("NodeType is not existed")
		return
	}
	// 删除过期数据
	for {
		time.Sleep(5 * time.Second)
		serverMap.Range(func(key, value interface{}) bool {
			lastTime := value.(time.Time)
			if time.Since(lastTime) > 10*time.Second {
				serverMap.Delete(key)
				addr, _ := key.(string)
				lbm.RemoveNode(t, addr)
			}
			return true
		})
	}
}

func GetDataServers() []string {
	return GetServers(common.DATANODE)
}

func GetApiServers() []string {
	return GetServers(common.APINODE)
}

func GetMasterServers() []string {
	return GetServers(common.MASTER)
}

func GetUserServers() []string {
	return GetServers(common.USERNODE)
}

func GetServers(t common.NodeType) []string {
	dataServers, existed := NodeMap[t]
	if !existed {
		log.Println("NodeType is not existed")
		return nil
	}
	ds := make([]string, 0)
	dataServers.Range(func(key, _ interface{}) bool {
		ds = append(ds, key.(string))
		return true
	})
	return ds
}

func StartHeartbeat() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()

	for {
		q.Publish(common.EXCHANGE_MASTER, types.HeartBeatMessage{
			Addr:           config.GetLocalAddr(),
			Weight:         config.GetNodeWeight(),
			NodeType:       common.MASTER,
			ResourceStatus: resource.NewResourceStatus(),
		})
		time.Sleep(common.HEATBEAT_INTERVAL * time.Second)
	}
}
