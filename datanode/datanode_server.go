package main

import (
	"DisHub/config"
	"DisHub/datanode/heartbeat"
	"DisHub/datanode/locate"
	"DisHub/datanode/objects"
	"DisHub/datanode/temp"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configFilePath string
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/data_server1.json", "Path to the config file")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}
}
func main() {
	// 加载已有文件
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	r := gin.Default()
	objects.Handler("/objects", r)
	temp.Handler("/temp", r)
	r.Run(config.GetLocalAddr())
}
