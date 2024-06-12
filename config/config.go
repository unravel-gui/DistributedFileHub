package config

import (
	"DisHub/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

type Config struct {
	IP           string `json:"ip"`
	PORT         int    `json:"port"`
	ENDPOINT     string
	MYSQLADDR    string `json:"mysqlAddr"`
	RABBITMQADDR string `json:"rabbitMQAddr"`
	BASEDIR      string `json:"basedir"`
	filePath     string
	Weight       int `json:"weight"`
	LBStrategy   int `json:"lbStrategy"`
	LBRetries    int `json:"lbRetries"`
}

var DefaultCfg Config

func (cfg *Config) LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("open config file err :%v", err)
		return err
	}
	defer file.Close()

	// json
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(cfg); err != nil {
		log.Printf("decode config file err: %v", err)
		return err
	}
	cfg.filePath = filePath
	cfg.ENDPOINT = cfg.IP + ":" + strconv.Itoa(cfg.PORT)
	utils.EnsureFileExists(cfg.BASEDIR)
	utils.EnsureFileExists(cfg.BASEDIR + "/objects")
	utils.EnsureFileExists(cfg.BASEDIR + "/temp")
	return nil
}

func GetRabbitMQAddr() string {
	addr := os.Getenv("RABBITMQ_SERVER")
	if addr == "" {
		addr = DefaultCfg.RABBITMQADDR
	}
	return addr
}

func GetLocalAddr() string {
	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		addr = DefaultCfg.ENDPOINT
	}
	return addr
}

func GetPort() int {
	port := DefaultCfg.PORT
	return port
}

func GetMySQLAddr() string {
	addr := os.Getenv("MYSQL_ADDRESS")
	if addr == "" {
		addr = DefaultCfg.MYSQLADDR
	}
	fmt.Println(addr)
	return addr
}

func GetBasePath() string {
	return DefaultCfg.BASEDIR
}

func GetFilePath(fileName string) string {
	return path.Join(DefaultCfg.BASEDIR, "/objects/"+fileName)
}

func GetLoadBalancerConfig() (int, int) {
	storage, retry := 0, 3
	if DefaultCfg.LBStrategy != 0 {
		storage = DefaultCfg.LBStrategy
	}
	if DefaultCfg.LBRetries != 3 {
		retry = DefaultCfg.LBRetries
	}
	return storage, retry
}

func GetNodeWeight() int {
	weight := 1
	if DefaultCfg.Weight > 1 {
		weight = DefaultCfg.Weight
	}
	return weight
}
