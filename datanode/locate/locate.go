package locate

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/config"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	return ok
}

func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	q := rabbitmq.New(config.GetRabbitMQAddr())
	defer q.Close()
	q.Bind(common.EXCHANGE_DATA)
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			fmt.Println("locate err", e)
		}
		exist := Locate(hash)
		fmt.Println("locate hash res:%v", exist)
		if exist {
			q.Send(msg.ReplyTo, config.GetLocalAddr())
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(config.GetBasePath() + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
