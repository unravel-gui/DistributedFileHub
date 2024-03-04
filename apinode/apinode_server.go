package main

import (
	"DisHub/apinode/heartbeat"
	"DisHub/apinode/locate"
	"DisHub/apinode/objects"
	"DisHub/config"
	service "DisHub/service/file_meta"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var listenAddress string

var (
	configFilePath string
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/api_server.json", "Path to the config file")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}
	service.G_fileMeta.Load()
}

func main() {
	go heartbeat.ListenHeartbeat()
	r := gin.Default()
	objects.Handler("/objects", r)
	locate.Handler("/locate", r)
	r.Run(config.GetLocalAddr())
}
