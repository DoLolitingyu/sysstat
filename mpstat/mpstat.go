package mpstat

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStats struct {
	CPU    string
	Usr    float64
	Nice   float64
	Sys    float64
	Iowait float64
	Irq    float64
	Soft   float64
	Steal  float64
	Guest  float64
	Gnice  float64
	Idle   float64
}

func GetData() ([]CPUStats,error){
	prevStats, err := GetCPUStats()
	if err != nil {
		return []CPUStats{}, err
	}

	time.Sleep(1 * time.Second) // should wait second

	currStats, err := GetCPUStats()
	if err != nil {
		return []CPUStats{}, err
	}

	stats := CalculateCPUUsage(prevStats, currStats)
	return stats, nil
}

func GetCPUStats() ([]CPUStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var stats []CPUStats

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) < 11 {
				continue
			}

			cpu := fields[0]
			stat := CPUStats{
				CPU:    cpu,
				Usr:    parseFloat(fields[1]),
				Nice:   parseFloat(fields[2]),
				Sys:    parseFloat(fields[3]),
				Idle:   parseFloat(fields[4]),
				Iowait: parseFloat(fields[5]),
				Irq:    parseFloat(fields[6]),
				Soft:   parseFloat(fields[7]),
				Steal:  parseFloat(fields[8]),
				Guest:  parseFloat(fields[9]),
				Gnice:  parseFloat(fields[10]),
			}
			stats = append(stats, stat)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func CalculateCPUUsage(prevStats, currStats []CPUStats) []CPUStats {
	var stats []CPUStats

	for i := 0; i < len(prevStats); i++ {
		prev := prevStats[i]
		curr := currStats[i]

		totalPrev := prev.Usr + prev.Nice + prev.Sys + prev.Iowait + prev.Irq + prev.Soft + prev.Steal + prev.Guest + prev.Gnice + prev.Idle
		totalCurr := curr.Usr + curr.Nice + curr.Sys + curr.Iowait + curr.Irq + curr.Soft + curr.Steal + curr.Guest + curr.Gnice + curr.Idle

		// calculate percent
		stat := CPUStats{
			CPU:    curr.CPU,
			Usr:    (curr.Usr - prev.Usr) * 100 / (totalCurr - totalPrev),
			Nice:   (curr.Nice - prev.Nice) * 100 / (totalCurr - totalPrev),
			Sys:    (curr.Sys - prev.Sys) * 100 / (totalCurr - totalPrev),
			Iowait: (curr.Iowait - prev.Iowait) * 100 / (totalCurr - totalPrev),
			Irq:    (curr.Irq - prev.Irq) * 100 / (totalCurr - totalPrev),
			Soft:   (curr.Soft - prev.Soft) * 100 / (totalCurr - totalPrev),
			Steal:  (curr.Steal - prev.Steal) * 100 / (totalCurr - totalPrev),
			Guest:  (curr.Guest - prev.Guest) * 100 / (totalCurr - totalPrev),
			Gnice:  (curr.Gnice - prev.Gnice) * 100 / (totalCurr - totalPrev),
			Idle:   (curr.Idle - prev.Idle) * 100 / (totalCurr - totalPrev),
		}

		stats = append(stats, stat)
	}

	return stats
}