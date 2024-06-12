package repository

type OSSMetadata struct {
	Hash  string `gorm:"primary_key;autoIncrement;column:hash"`
	Size  int64  `gorm:"column:size"`
	IsDel int    `gorm:"column:is_del"`
}

func (OSSMetadata) TableName() string {
	return "oss_meta" // 设置表名
}
