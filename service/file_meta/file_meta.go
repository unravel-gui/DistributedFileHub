package service

type FileMetadata struct {
	Filename string `gorm:"column:name"`
	Size     int64  `gorm:"column:size"`
	Hash     string `gorm:"column:hash"`
	IsDel    int    `gorm:"column:is_del"`
}

func (FileMetadata) TableName() string {
	return "file_meta" // 设置表名
}
