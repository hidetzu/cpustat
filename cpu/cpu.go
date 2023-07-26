package cpu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Get cpu statistics
func Get() (*Stats, error) {
	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return collectCPUStats(file)
}

// Stats represents cpu statistics for linux
type Stats struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
	Total     uint64
	CPUCount  int
	StatCount int

	UserPercent   float64
	NicePercent   float64
	SystemPercent float64
	IdlePercent   float64
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func collectCPUStats(out io.Reader) (*Stats, error) {
	scanner := bufio.NewScanner(out)
	var cpu Stats

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"nice", &cpu.Nice},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"iowait", &cpu.Iowait},
		{"irq", &cpu.Irq},
		{"softirq", &cpu.Softirq},
		{"steal", &cpu.Steal},
		{"guest", &cpu.Guest},
		{"guest_nice", &cpu.GuestNice},
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to scan /proc/stat")
	}

	valStrs := strings.Fields(scanner.Text())[1:]
	cpu.StatCount = len(valStrs)
	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s from /proc/stat", cpuStats[i].name)
		}
		*cpuStats[i].ptr = val
		cpu.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cpu.Total -= cpu.Guest
	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	cpu.Total -= cpu.GuestNice

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			cpu.CPUCount++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for /proc/stat: %s", err)
	}

	cpu.UserPercent = float64(cpu.User) / float64(cpu.Total) * 100
	cpu.NicePercent = float64(cpu.Nice) / float64(cpu.Total) * 100
	cpu.SystemPercent = float64(cpu.System) / float64(cpu.Total) * 100
	cpu.IdlePercent = float64(cpu.Idle) / float64(cpu.Total) * 100

	return &cpu, nil
}
