package repository

import (
	"DisHub/common"
	"DisHub/common/resource"
	"time"
)

type ResourceUsage struct {
	Rid      int             `gorm:"primary_key;autoIncrement" json:"rid"`
	Addr     string          `json:"addr"`
	NodeType common.NodeType `json:"node_type"`
	resource.ResourceStatus
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (ResourceUsage) TableName() string {
	return "resource_usage"
}
