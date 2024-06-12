package repository

import (
	"DisHub/common/db"
	"DisHub/common/utils"
	"DisHub/config"
	"gorm.io/gorm"
	"log"
	"time"
)

type ResourceUsageRepository struct {
	db *gorm.DB
}

func NewResourceUsageRepository(mysqlAddr string) *ResourceUsageRepository {
	if mysqlAddr == "" {
		mysqlAddr = config.GetLocalAddr()
	}
	d, err := db.NewConnect(mysqlAddr)
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	d.AutoMigrate(&ResourceUsage{})
	return &ResourceUsageRepository{
		db: d,
	}
}

func (rur *ResourceUsageRepository) Close() {
	sqlDB, err := rur.db.DB()
	if err != nil {
		log.Printf("close Mysql connect err:%v\n", err)
		return
	}
	sqlDB.Close()
}

func (rur *ResourceUsageRepository) Insert(ru *ResourceUsage) error {
	result := rur.db.Create(ru)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (rur *ResourceUsageRepository) GetResource(addr string, interval time.Duration) ([]*ResourceUsage, error) {
	// 计算起始时间
	startTime := utils.GetNow().Add(-interval)
	// 查询数据库，获取指定时间段内的资源使用记录
	var resourceUsageList []*ResourceUsage
	if err := rur.db.
		Select("create_time, cpu_usage, memory_usage, disk_usage").
		Where("addr=? AND create_time >= ?", addr, startTime).
		Order("create_time DESC").
		Limit(10).
		Find(&resourceUsageList).Error; err != nil {
		return nil, err
	}
	return resourceUsageList, nil
}

func (rur *ResourceUsageRepository) GetResourceLastest() []ResourceUsage {
	var resourceUsages []ResourceUsage
	// 执行 SQL 查询
	query := `
        SELECT a.*
		FROM resource_usage AS a
		INNER JOIN (
			SELECT addr, MAX(create_time) AS max_create_time
			FROM resource_usage
			WHERE create_time >= DATE_SUB(NOW(), INTERVAL 10 MINUTE)
			GROUP BY addr
		) AS b ON a.addr = b.addr AND a.create_time = b.max_create_time
		WHERE a.create_time >= DATE_SUB(NOW(), INTERVAL 10 MINUTE)
    `
	rur.db.Raw(query).Scan(&resourceUsages)

	return resourceUsages
}
