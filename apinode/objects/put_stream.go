package objects

import (
	"DisHub/apinode/heartbeat"
	"DisHub/common/objectstream"
	"fmt"
	"log"
)

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	log.Println(server)
	return objectstream.NewTempPutStream(server, hash, size)
}
