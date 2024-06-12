package repository

import "time"

type Task struct {
	Tid        int       `gorm:"primary_key;autoIncrement" json:"tid"`
	Uid        int       `gorm:"column:uid" json:"uid"`   // 文件对应的用户id
	Hash       string    `gorm:"column:hash" json:"hash"` // 文件的hash值
	Size       int64     `gorm:"column:size" json:"size"` // 文件的大小
	Type       string    `gorm:"column:type" json:"type"` // 文件类型，可以是目录"common.FOLDER"
	Path       string    `gorm:"column:path" json:"path"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
}
