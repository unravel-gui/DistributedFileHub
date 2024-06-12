package resource

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type ResourceStatus struct {
	CPUUsage        float64 `gorm:"column:cpu_usage" json:"cpu_usage"`                 // CPU利用率
	CPUCurrentUsage float64 `gorm:"column:cpu_current_usage" json:"cpu_current_usage"` // CPU当前使用量
	CPUMaxUsage     float64 `gorm:"column:cpu_max_usage" json:"cpu_max_usage"`         // CPU总量

	MemoryUsage        float64 `gorm:"column:memory_usage" json:"memory_usage"`                 // 内存利用率
	MemoryCurrentUsage uint64  `gorm:"column:memory_current_usage" json:"memory_current_usage"` // 内存当前使用量
	MemoryTotal        uint64  `gorm:"column:memory_total" json:"memory_total"`                 // 内存总量

	DiskUsage        float64 `gorm:"column:disk_usage" json:"disk_usage"`                 // 磁盘利用率
	DiskCurrentUsage uint64  `gorm:"column:disk_current_usage" json:"disk_current_usage"` // 磁盘当前使用量
	DiskTotal        uint64  `gorm:"column:disk_total" json:"disk_total"`                 // 磁盘总量

	NetworkUsage float64 `gorm:"column:network_usage" json:"network_usage"` // 网络利用率
	NetworkSent  uint64  `gorm:"column:network_sent" json:"network_sent"`   // 网络发送量
	NetworkRecv  uint64  `gorm:"column:network_recv" json:"network_recv"`   // 网络接收量
}

func NewResourceStatus() ResourceStatus {
	// 获取 CPU 利用率
	cpuPercent, _ := cpu.Percent(0, false)
	// 获取 CPU 核心数
	cpuCount := float64(len(cpuPercent))

	// 获取内存信息
	memInfo, _ := mem.VirtualMemory()

	// 获取磁盘信息
	diskUsage, _ := disk.Usage("/")

	// 获取网络信息
	netStats, _ := net.IOCounters(false)
	var networkSent, networkRecv uint64
	for _, stats := range netStats {
		networkSent += stats.BytesSent
		networkRecv += stats.BytesRecv
	}

	return ResourceStatus{
		CPUUsage:        cpuPercent[0],
		CPUCurrentUsage: cpuPercent[0] * cpuCount / 100,
		CPUMaxUsage:     cpuCount,

		MemoryUsage:        memInfo.UsedPercent,
		MemoryCurrentUsage: memInfo.Used,
		MemoryTotal:        memInfo.Total,

		DiskUsage:        diskUsage.UsedPercent,
		DiskCurrentUsage: diskUsage.Used,
		DiskTotal:        diskUsage.Total,

		NetworkUsage: float64(networkSent + networkRecv),
		NetworkSent:  networkSent,
		NetworkRecv:  networkRecv,
	}
}
