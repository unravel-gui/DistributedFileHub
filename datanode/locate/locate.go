package locate

import (
	"DisHub/common"
	"DisHub/common/rabbitmq"
	"DisHub/common/types"
	"DisHub/config"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
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
		id := Locate(hash)
		log.Println(id)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{
				Addr: config.GetLocalAddr(),
				Id:   id,
			})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(config.GetBasePath() + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			fmt.Println("collected filename=", file)
			continue
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			fmt.Println("collect objects err:", e)
		}
		objects[hash] = id
	}
}
