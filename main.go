package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

const cpuThreshold = 80.0
const memThreshold = 10.0

func main() {
	for {
		checkSystemUsage()
		time.Sleep(5 * time.Second) // Check every 10 seconds
	}
}

func checkSystemUsage() {
	cpuUsage, _ := cpu.Percent(0, false)
	memUsage, _ := mem.VirtualMemory()

	if cpuUsage[0] > cpuThreshold || memUsage.UsedPercent > memThreshold {
		generateReport()
	}
}

func generateReport() {
	fmt.Println("Generating report...")

	// Get process details
	processes, _ := process.Processes()
	for _, proc := range processes {
		cpuPercent, _ := proc.CPUPercent()
		memInfo, _ := proc.MemoryInfo()
		name, _ := proc.Name()
		fmt.Printf("Process: %s, CPU: %.2f%%, Memory: %v bytes\n", name, cpuPercent, memInfo.RSS)
	}

	// Get network connections
	connections, _ := net.Connections("all")
	for _, conn := range connections {
		fmt.Printf("Network connection: %s:%d -> %s:%d, Status: %s\n", conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status)
	}

	saveReportToFile()
}

func saveReportToFile() {
	file, err := os.Create("report.txt")
	if err != nil {
		fmt.Println("Error creating report file:", err)
		return
	}
	defer file.Close()

	// Add system usage details
	cpuUsage, _ := cpu.Percent(0, false)
	memUsage, _ := mem.VirtualMemory()

	file.WriteString(fmt.Sprintf("CPU Usage: %.2f%%\n", cpuUsage[0]))
	file.WriteString(fmt.Sprintf("Memory Usage: %.2f%%\n\n", memUsage.UsedPercent))

	// Add process details
	processes, _ := process.Processes()
	for _, proc := range processes {
		cpuPercent, _ := proc.CPUPercent()
		memInfo, _ := proc.MemoryInfo()
		name, _ := proc.Name()
		file.WriteString(fmt.Sprintf("Process: %s, CPU: %.2f%%, Memory: %v bytes\n", name, cpuPercent, memInfo.RSS))
	}

	// Add network connections
	connections, _ := net.Connections("all")
	for _, conn := range connections {
		file.WriteString(fmt.Sprintf("Network connection: %s:%d -> %s:%d, Status: %s\n", conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status))
	}
}
