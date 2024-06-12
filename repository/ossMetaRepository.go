package repository

import (
	"DisHub/common/db"
	"DisHub/config"
	"errors"
	"gorm.io/gorm"
	"log"
)

type OSSMetaRepository struct {
	db *gorm.DB
}

func NewOSSMetaRepository(mysqlAddr string) *OSSMetaRepository {
	if mysqlAddr == "" {
		mysqlAddr = config.GetLocalAddr()
	}
	d, err := db.NewConnect(mysqlAddr)
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	d.AutoMigrate(&OSSMetadata{})
	return &OSSMetaRepository{
		db: d,
	}
}

func (oms *OSSMetaRepository) Close() {
	sqlDB, err := oms.db.DB()
	if err != nil {
		log.Printf("close Mysql connect err:%v\n", err)
		return
	}
	sqlDB.Close()
}

func (oms *OSSMetaRepository) Insert(om *OSSMetadata) error {
	result := oms.db.Create(om)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (oms *OSSMetaRepository) GetMetaData(hash string) (*OSSMetadata, error) {
	var fm OSSMetadata
	result := oms.db.Where("hash = ? and is_del=0", hash).First(&fm)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &fm, nil
}

func (oms *OSSMetaRepository) DelMetaData(hash string) error {
	result := oms.db.Model(&OSSMetadata{}).Where("hash = ? AND is_del = ?", hash, 0).Update("is_del", 1)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	return nil
}
