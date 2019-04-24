package health

import (
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"math"
)

const (
	ErrorDomain         = "HEALTH_COLLECTOR"
	DiskUsageErrCode    = "DISK_USAGE_ERROR"
	NetworkUsageErrCode = "NETWORK_USAGE_ERROR"
	CPUUsageErrorCode   = "CPU_USAGE_ERROR"
)

// Server stats dto
type ServerStats struct {
	NetworkStats float64
	DiskUsage    float64
	CPUStats     float64
}

// Server load stats collector service
type LoadStatsCollector struct {
	errHandler *kerr.Handler
}

// Constructor
func NewLoadStatsCollector(errHandler *kerr.Handler) *LoadStatsCollector {
	collector := new(LoadStatsCollector)
	collector.errHandler = errHandler

	return collector
}

// Collect utilization parameters
func (hs *LoadStatsCollector) CollectHealthData() ServerStats {
	stats := ServerStats{}

	diskStat, err := disk.Usage("/")
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlWarning, ErrorDomain, DiskUsageErrCode, err, nil)

		hs.errHandler.Handle(e)
	}
	stats.DiskUsage = math.Round(diskStat.UsedPercent)

	netStat := 0.0
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlWarning, ErrorDomain, NetworkUsageErrCode, err, nil)

		hs.errHandler.Handle(e)
	}
	stats.NetworkStats = netStat

	cpuStats, err := cpu.Percent(10, false)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlWarning, ErrorDomain, CPUUsageErrorCode, err, nil)

		hs.errHandler.Handle(e)
	}
	stats.CPUStats = cpuStats[0]

	return stats
}
