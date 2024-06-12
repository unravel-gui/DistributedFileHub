package heartbeat

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/common/resource"
	"DisHub/common/stub"
	"DisHub/common/types"
	"DisHub/config"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// 保存提供服务的datanode节点
var masterServers []string
var mutex sync.Mutex

func getMasterAddr() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	q.Publish(common.EXCHANGE_LOCATE_MASTER, config.GetLocalAddr())
	c := q.Consume()
	masterServers = nil
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	for msg := range c {
		var server string
		json.Unmarshal(msg.Body, &server)
		if server == "" {
			continue
		}
		masterServers = append(masterServers, server)
	}
}

func StartListenMaster() {
	mutex.Lock()
	getMasterAddr()
	mutex.Unlock()
	for {
		mutex.Lock()
		getMasterAddr()
		mutex.Unlock()
		time.Sleep(15 * time.Second)
	}
}

func GetMasterServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	dd := make([]string, len(masterServers))
	copy(dd, masterServers)
	return dd
}

func GetDataServers() ([]string, error) {
	return getServers(common.DATANODE)
}

func GetUserServers() ([]string, error) {
	return getServers(common.USERNODE)
}

func getServers(t common.NodeType) ([]string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	servers := masterServers
	if len(servers) == 0 {
		return nil, errors.New("no enough master")
	}
	common.G_Random.Shuffle(len(servers), func(i, j int) {
		servers[i], servers[j] = servers[j], servers[i]
	})
	ds := make([]string, 0)
	for _, server := range servers {
		switch t {
		case common.DATANODE:
			ds = stub.GetDataServers(server)
		case common.USERNODE:
			ds = stub.GetUserServers(server)
		}
		if len(ds) != 0 {
			break
		}
	}
	return ds, nil
}

func GetBalanceUserServer() (string, error) {
	return getBalanceServer(common.USERNODE)
}

func getBalanceServer(t common.NodeType) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	servers := masterServers
	if len(servers) == 0 {
		return "", errors.New("no enough master")
	}
	common.G_Random.Shuffle(len(servers), func(i, j int) {
		servers[i], servers[j] = servers[j], servers[i]
	})
	ds := ""
	for _, server := range masterServers {
		switch t {
		case common.USERNODE:
			ds = stub.GetGetLoadBalanceUserServer(server)
		}
		if ds != "" {
			break
		}
	}

	return ds, nil
}

func StartHeartbeat() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()

	for {
		q.Publish(common.EXCHANGE_MASTER, types.HeartBeatMessage{
			Addr:           config.GetLocalAddr(),
			Weight:         config.GetNodeWeight(),
			NodeType:       common.APINODE,
			ResourceStatus: resource.NewResourceStatus(),
		})
		time.Sleep(common.HEATBEAT_INTERVAL * time.Second)
	}
}
