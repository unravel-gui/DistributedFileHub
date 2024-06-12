package common

import (
	"math/rand"
	"strings"
	"time"
)

const (
	EXCHANGE_MASTER        = "masterServers"
	EXCHANGE_LOCATE_MASTER = "locateMaster"
	EXCHANGE_API           = "apiServers"
	EXCHANGE_DATA          = "dataServers"
)

type NodeType int64

const (
	MASTER NodeType = iota
	DATANODE
	APINODE
	USERNODE
)

func StringToNodeType(nt string) NodeType {
	nt = strings.ToLower(nt)
	switch nt {
	case "master":
		return MASTER
	case "datanode":
		return DATANODE
	case "apinode":
		return APINODE
	case "usernode":
		return USERNODE
	default:
		return -1 // 表示未知的节点类型
	}
}

func (nt NodeType) ToString() string {
	switch nt {
	case MASTER:
		return "MASTER"
	case DATANODE:
		return "DATANODE"
	case APINODE:
		return "APINODE"
	case USERNODE:
		return "USERNODE"
	}
	return "unknown"
}

const (
	RPC_TOKEN_KEY   = "magic_token"
	RPC_TOKEN_VALVE = "rpc_token"
)

const (
	FOLDER     = "Folder"
	ROOTFOLDER = -1
)

const (
	TEST_MYSQL_ADDR = "ossfile:ossfile@tcp(localhost:3306)/oss_fileMeta?parseTime=true"
	TEST_UID        = -1
)

const (
	JWT_KEY   = "jwt_key"
	JWT_TOKEN = "Jwt_Token"
)

const (
	REGULAR = "regular"
	VIDEOS  = "videos"
	IMAGES  = "images"
)

const HEATBEAT_INTERVAL = 15

var G_Random = rand.New(rand.NewSource(time.Now().UnixNano()))
