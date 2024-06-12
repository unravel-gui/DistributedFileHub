package main

import (
	"DisHub/config"
	"DisHub/middleware"
	"DisHub/service"
	"DisHub/usernode/fileMeta"
	"DisHub/usernode/heartbeat"
	"DisHub/usernode/user"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configFilePath string
	enableFlag     bool
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/usernode1.json", "Path to the config file")
	flag.BoolVar(&enableFlag, "c", false, "terminate")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}
	service.G_User.Load()
	service.G_FileMeta.Load()
}

func main() {
	//pidFileName := fmt.Sprintf("%s/usernode_%d.pid", config.GetBasePath(), config.GetPort())
	//if enableFlag {
	//	err := utils.KillPID(pidFileName)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	return
	//} else {
	//	if utils.FileExists(pidFileName) {
	//		log.Fatalf("locked by %s", pidFileName)
	//	}
	//	err := utils.SavePID(pidFileName)
	//	if err != nil {
	//		log.Fatalf("create pid file err: %v", err)
	//	}
	//}

	// 加载已有文件
	go heartbeat.StartHeartbeat()
	r := gin.Default()
	//r.Use(middleware.CorsMiddleware())
	user.HandlerWithoutCheck("/", r)
	r.Use(middleware.CheckJWTToken())
	user.Handler("/user", r)
	fileMeta.Handler("/fileMeta", r)
	r.Run(config.GetLocalAddr())
}
