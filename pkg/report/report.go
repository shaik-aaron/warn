package report

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func GnerateReport() {
	fmt.Println("Generating report...")

	// Create or open the report file
	file, err := os.Create("report.txt")
	if err != nil {
		fmt.Println("Error creating report file:", err)
		return
	}
	defer file.Close()

	// Add system usage details to file
	if err := writeSystemUsage(file); err != nil {
		fmt.Println("Error adding system usage details:", err)
	}

	// Collect and print process details
	if err := writeProcessDetails(file); err != nil {
		fmt.Println("Error writing process details:", err)
	}

	// Collect and print network connections
	if err := writeNetworkConnections(file); err != nil {
		fmt.Println("Error writing network connections:", err)
	}
}

func writeSystemUsage(file *os.File) error {
	// Add CPU usage details
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Error retrieving CPU usage:", err)
		return err
	}
	file.WriteString(fmt.Sprintf("CPU Usage: %.2f%%\n", cpuUsage[0]))

	// Add memory usage details
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error retrieving memory usage:", err)
		return err
	}
	file.WriteString(fmt.Sprintf("Memory Usage: %.2f%%\n\n", memUsage.UsedPercent))

	return nil
}

func writeProcessDetails(file *os.File) error {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error retrieving process details:", err)
		return err
	}

	processMap := make(map[string]int)
	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}
		processMap[name]++
	}

	for name, count := range processMap {
		fmt.Printf("Process: %s, Count: %d\n", name, count)
		file.WriteString(fmt.Sprintf("Process: %s, Count: %d\n", name, count))
		for _, proc := range processes {
			procName, err := proc.Name()
			if err != nil || procName != name {
				continue
			}
			cpuPercent, err := proc.CPUPercent()
			if err != nil {
				continue
			}
			memInfo, err := proc.MemoryInfo()
			if err != nil {
				continue
			}
			fmt.Printf("\tCPU: %.2f%%, Memory: %v bytes\n", cpuPercent, memInfo.RSS)
			file.WriteString(fmt.Sprintf("\tCPU: %.2f%%, Memory: %v bytes\n", cpuPercent, memInfo.RSS))
		}
	}
	return nil
}

func writeNetworkConnections(file *os.File) error {
	connections, err := net.Connections("all")
	if err != nil {
		fmt.Println("Error retrieving network connections:", err)
		return err
	}

	for _, conn := range connections {
		fmt.Printf("Network connection: %s:%d -> %s:%d, Status: %s\n", conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status)
		file.WriteString(fmt.Sprintf("Network connection: %s:%d -> %s:%d, Status: %s\n", conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status))
	}
	return nil
}