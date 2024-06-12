package main

import (
	"DisHub/apinode/avator"
	"DisHub/apinode/heartbeat"
	"DisHub/apinode/locate"
	"DisHub/apinode/objects"
	"DisHub/apinode/proxy"
	"DisHub/apinode/temp"
	"DisHub/config"
	"DisHub/middleware"
	"DisHub/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var listenAddress string

var (
	configFilePath string
	enableFlag     bool
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/api_server1.json", "Path to the config file")
	flag.BoolVar(&enableFlag, "c", false, "terminate")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}
	service.G_OssMeta.Load()
}

func main() {
	//pidFileName := fmt.Sprintf("%s/apinode_%d.pid", config.GetBasePath(), config.GetPort())
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

	go heartbeat.StartHeartbeat()
	go heartbeat.StartListenMaster()

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	proxy.UserHandlerWithOutCheck("/", r)
	// 头像
	avator.Handler("/avator", r)
	// 加入JWT验证
	r.Use(middleware.CheckJWTToken())
	objects.Handler("/objects", r)
	temp.Handler("/temp", r)
	locate.Handler("/locate", r)
	proxy.UserHandler("/user", r)
	proxy.UserHandler("/fileMeta", r)
	proxy.MasterHandler("/resourceUsage", r)
	proxy.MasterHandler("/lb", r)
	r.Run(config.GetLocalAddr())
}
