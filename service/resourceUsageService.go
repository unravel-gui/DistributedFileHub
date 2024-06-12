package service

import (
	"DisHub/common/types"
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/repository"
	"time"
)

type ResourceUsageService struct {
	r *repository.ResourceUsageRepository
}

var G_ResourceUsage ResourceUsageService

func NewResourceUsageService(mysqlAddr string) *ResourceUsageService {
	if mysqlAddr == "" {
		mysqlAddr = config.GetMySQLAddr()
	}
	r := repository.NewResourceUsageRepository(mysqlAddr)
	return &ResourceUsageService{
		r: r,
	}
}

func (rur *ResourceUsageService) Load() {
	mysqlAddr := config.GetMySQLAddr()
	r := repository.NewResourceUsageRepository(mysqlAddr)
	rur.r = r
}

func (rur *ResourceUsageService) Close() {
	rur.r.Close()
}

func (rur *ResourceUsageService) InsertOne(info types.HeartBeatMessage) error {
	rs := &repository.ResourceUsage{
		Addr:           info.Addr,
		NodeType:       info.NodeType,
		ResourceStatus: info.ResourceStatus,
		CreateTime:     utils.GetNow(),
	}
	return rur.r.Insert(rs)
}

func (rur *ResourceUsageService) GetNodeResourceUsage(addr string) ([]*repository.ResourceUsage, error) {
	interval := 1 * time.Hour
	return rur.r.GetResource(addr, interval)
}

func (rur *ResourceUsageService) GetResourceLastest() map[string][]repository.ResourceUsage {
	rss := rur.r.GetResourceLastest()
	grs := make(map[string][]repository.ResourceUsage)
	for _, rs := range rss {
		grs[rs.NodeType.ToString()] = append(grs[rs.NodeType.ToString()], rs)
	}
	return grs
}
