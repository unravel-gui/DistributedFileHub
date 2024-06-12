package repository

import (
	"DisHub/common/utils"
	"time"
)

type FileMeta struct {
	Fid         int       `gorm:"primary_key;autoIncrement" json:"fid"`
	Uid         int       `gorm:"column:uid" json:"uid"`                   // 文件对应的用户id
	Dir         int       `gorm:"column:dir;" json:"dir"`                  // 文件位于的目录
	Hash        string    `gorm:"column:hash" json:"hash"`                 // 文件的hash值
	Name        string    `gorm:"column:name;" json:"name"`                //文件名称，建立联合唯一索引
	Size        int64     `gorm:"column:size" json:"size"`                 // 文件的大小
	ContentType string    `gorm:"column:content_type" json:"content_type"` // 文件类型，可以是目录"common.FOLDER"
	IsDel       bool      `gorm:"column:is_del;" json:"_"`                 // 上传时间或创建时间
	UploadTime  time.Time `gorm:"column:upload_time" json:"upload_time"`   // 上传时间或创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`   // 上传时间或创建时间
}

func (FileMeta) TableName() string {
	return "file_meta" // 设置表名
}

func (f *FileMeta) UpdateFileMeta(newFileMeta *FileMeta) {
	f.Dir = newFileMeta.Dir
	f.Name = newFileMeta.Name
	f.UpdateTime = utils.GetNow()
}

func (f *FileMeta) CommitFileMeta() bool {

	return true
}
