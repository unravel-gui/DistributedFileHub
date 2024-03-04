package service

import (
	"DisHub/common/db"
	"DisHub/config"
	"errors"
	"gorm.io/gorm"
	"log"
)

type FileMetaService struct {
	db *gorm.DB
}

var G_fileMeta FileMetaService

func NewFileMetaService() *FileMetaService {
	d, err := db.NewConnect(config.GetMySQLAddr())
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	return &FileMetaService{
		db: d,
	}
}

func (fms *FileMetaService) Load() {
	d, err := db.NewConnect(config.GetMySQLAddr())
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	fms.db = d
}

func (fms *FileMetaService) Close() {
	sqlDB, err := fms.db.DB()
	if err != nil {
		log.Printf("close Mysql connect err:%v\n", err)
		return
	}
	sqlDB.Close()
}

func (fms *FileMetaService) PutMetaData(name string, size int64, hash string) error {
	fm := FileMetadata{
		Filename: name,
		Size:     size,
		Hash:     hash,
	}
	result := fms.db.Create(fm)
	if result.Error != nil {
		log.Println("InsertOne error: ", result.Error)
		return result.Error
	}
	return nil
}

func (fms *FileMetaService) GetLatest(name string) (*FileMetadata, error) {
	var fm FileMetadata
	result := fms.db.Where("name = ? and is_del=0", name).First(&fm)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &fm, nil
}

func (fms *FileMetaService) GetMetaData(hash string) (*FileMetadata, error) {
	var fm FileMetadata
	result := fms.db.Where("hash = ? and is_del=0", hash).First(&fm)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &fm, nil
}

func (fms *FileMetaService) DelMetaData(hash string) error {
	result := fms.db.Model(&FileMetadata{}).Where("hash = ? AND is_del = ?", hash, 0).Update("is_del", 1)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	return nil
}
