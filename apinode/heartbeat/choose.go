package heartbeat

import (
	"log"
	"math/rand"
)

func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		log.Println("no data server to choose")
		return ""
	}
	return ds[rand.Intn(n)]
}
