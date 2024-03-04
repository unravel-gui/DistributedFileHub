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
