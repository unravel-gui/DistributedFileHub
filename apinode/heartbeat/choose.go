package heartbeat

import (
	"log"
	"math/rand"
)

func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidate := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	// 过滤需要排除的节点
	servers, _ := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidate = append(candidate, s)
		}
	}
	length := len(candidate)
	if length < n {
		return
	}
	// 打乱顺序
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidate[p[i]])
	}
	return
}

func ChooseRandomDataServer() string {
	ds, _ := GetDataServers()
	n := len(ds)
	if n == 0 {
		log.Println("no data server to choose")
		return ""
	}
	return ds[rand.Intn(n)]
}
