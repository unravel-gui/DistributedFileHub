package main

import (
	"DisHub/config"
	"DisHub/loadbalancer"
	"DisHub/master/apinode"
	"DisHub/master/balancer"
	"DisHub/master/datanode"
	"DisHub/master/heartbeat"
	"DisHub/master/locate"
	"DisHub/master/master"
	"DisHub/master/resourceUsage"
	"DisHub/master/usernode"
	"DisHub/middleware"
	"DisHub/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configFilePath string
	enableFlag     bool
)

func init() {
	flag.StringVar(&configFilePath, "f", "config/master1.json", "Path to the config file")
	flag.BoolVar(&enableFlag, "c", false, "terminate")
	flag.Parse()
	err := config.DefaultCfg.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config:%v", err)
		return
	}
	service.G_ResourceUsage.Load()
	loadbalancer.G_LoadBalancerMap.Load()
}

func main() {
	//pidFileName := fmt.Sprintf("%s/master_%d.pid", config.GetBasePath(), config.GetPort())
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

	go heartbeat.ListenHeartbeat()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	r := gin.Default()

	balancer.Handler("/lb", r)
	r.Use(middleware.CheckJWTToken())
	datanode.Handler("/datanode", r)
	apinode.Handler("/apinode", r)
	usernode.Handler("/usernode", r)
	master.Handler("/master", r)
	resourceUsage.Handler("/resourceUsage", r)
	r.Run(config.GetLocalAddr())
}
