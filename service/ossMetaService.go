package service

import (
	"DisHub/config"
	"DisHub/repository"
)

type OSSMetaService struct {
	r *repository.OSSMetaRepository
}

var G_OssMeta OSSMetaService

func NewOSSMetaService(mysqlAddr string) *OSSMetaService {
	if mysqlAddr == "" {
		mysqlAddr = config.GetMySQLAddr()
	}
	r := repository.NewOSSMetaRepository(mysqlAddr)
	return &OSSMetaService{
		r: r,
	}
}

func (oms *OSSMetaService) Load() {
	mysqlAddr := config.GetMySQLAddr()
	r := repository.NewOSSMetaRepository(mysqlAddr)
	oms.r = r
}

func (oms *OSSMetaService) Close() {
	oms.r.Close()
}

func (oms *OSSMetaService) PutMetaData(hash string, size int64) error {
	om := &repository.OSSMetadata{
		Hash: hash,
		Size: size,
	}
	return oms.r.Insert(om)
}

func (oms *OSSMetaService) GetMetaData(hash string) (*repository.OSSMetadata, error) {
	return oms.r.GetMetaData(hash)
}

func (oms *OSSMetaService) DelMetaData(hash string) error {
	return oms.r.DelMetaData(hash)
}
