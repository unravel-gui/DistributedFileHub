package main

import (
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/datanode/heartbeat"
	"DisHub/datanode/locate"
	"DisHub/datanode/objects"
	"DisHub/datanode/temp"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configFilePath string
	enableFlag     bool
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/data_server1.json", "Path to the config file")
	flag.BoolVar(&enableFlag, "c", false, "terminate")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}

}

func main() {
	pidFileName := fmt.Sprintf("%s/datanode_%d.pid", config.GetBasePath(), config.GetPort())
	if enableFlag {
		err := utils.KillPID(pidFileName)
		if err != nil {
			log.Fatalln(err)
		}
		return
	} else {
		if utils.FileExists(pidFileName) {
			log.Fatalf("locked by %s", pidFileName)
		}
		err := utils.SavePID(pidFileName)
		if err != nil {
			log.Fatalf("create pid file err: %v", err)
		}
	}

	// 加载已有文件
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	r := gin.Default()
	objects.Handler("/objects", r)
	temp.Handler("/temp", r)
	r.Run(config.GetLocalAddr())
}
