package types

import (
	"DisHub/common"
	"DisHub/common/resource"
)

type LocateMessage struct {
	Addr string
	Id   int
}

type HeartBeatMessage struct {
	Addr           string
	Weight         int
	NodeType       common.NodeType
	ResourceStatus resource.ResourceStatus
}
