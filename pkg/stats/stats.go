package stats

import (
	"warn/pkg/report"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const cpuThreshold = 80.0
const memThreshold = 10.0

func CheckSystemUsage() {
	cpuUsage, _ := cpu.Percent(0, false)
	memUsage, _ := mem.VirtualMemory()

	if cpuUsage[0] > cpuThreshold || memUsage.UsedPercent > memThreshold {
		report.GnerateReport()
	}
}
