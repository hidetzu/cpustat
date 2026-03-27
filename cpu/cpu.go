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

// CoreStats holds raw /proc/stat tick counters for one CPU line.
type CoreStats struct {
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
}

// Stats represents a snapshot of /proc/stat.
type Stats struct {
	CPU       CoreStats
	CPUCount  int
	StatCount int

	UserPercent   float64
	NicePercent   float64
	SystemPercent float64
	IdlePercent   float64
}

// fieldNames maps index to field name for error messages.
var fieldNames = []string{
	"user", "nice", "system", "idle", "iowait",
	"irq", "softirq", "steal", "guest", "guest_nice",
}

// parseLine parses the numeric fields from a /proc/stat CPU line.
// fields should be the space-separated tokens after the "cpu"/"cpuN" prefix.
func parseLine(fields []string) (CoreStats, error) {
	var cs CoreStats
	ptrs := []*uint64{
		&cs.User, &cs.Nice, &cs.System, &cs.Idle, &cs.Iowait,
		&cs.Irq, &cs.Softirq, &cs.Steal, &cs.Guest, &cs.GuestNice,
	}

	for i, f := range fields {
		if i >= len(ptrs) {
			break
		}
		val, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			name := fieldNames[i]
			return CoreStats{}, fmt.Errorf("failed to parse %s: %w", name, err)
		}
		*ptrs[i] = val
		cs.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cs.Total -= cs.Guest
	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	cs.Total -= cs.GuestNice

	return cs, nil
}

// Get reads /proc/stat and returns a CPU statistics snapshot.
func Get() (*Stats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	return collectCPUStats(file)
}

func collectCPUStats(out io.Reader) (*Stats, error) {
	scanner := bufio.NewScanner(out)

	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to scan /proc/stat")
	}

	fields := strings.Fields(scanner.Text())[1:]
	cs, err := parseLine(fields)
	if err != nil {
		return nil, err
	}

	s := &Stats{
		CPU:       cs,
		StatCount: len(fields),
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			s.CPUCount++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for /proc/stat: %s", err)
	}

	if s.CPU.Total > 0 {
		s.UserPercent = float64(s.CPU.User) / float64(s.CPU.Total) * 100
		s.NicePercent = float64(s.CPU.Nice) / float64(s.CPU.Total) * 100
		s.SystemPercent = float64(s.CPU.System) / float64(s.CPU.Total) * 100
		s.IdlePercent = float64(s.CPU.Idle) / float64(s.CPU.Total) * 100
	}

	return s, nil
}

// Delta calculates the CPU tick deltas between two snapshots.
// prev should be the earlier snapshot and next the later one.
// Returns nil if the total delta is zero (no time has passed).
func Delta(prev, next *Stats) *Stats {
	d := &Stats{
		CPU: CoreStats{
			User:      next.CPU.User - prev.CPU.User,
			Nice:      next.CPU.Nice - prev.CPU.Nice,
			System:    next.CPU.System - prev.CPU.System,
			Idle:      next.CPU.Idle - prev.CPU.Idle,
			Iowait:    next.CPU.Iowait - prev.CPU.Iowait,
			Irq:       next.CPU.Irq - prev.CPU.Irq,
			Softirq:   next.CPU.Softirq - prev.CPU.Softirq,
			Steal:     next.CPU.Steal - prev.CPU.Steal,
			Guest:     next.CPU.Guest - prev.CPU.Guest,
			GuestNice: next.CPU.GuestNice - prev.CPU.GuestNice,
			Total:     next.CPU.Total - prev.CPU.Total,
		},
		CPUCount:  next.CPUCount,
		StatCount: next.StatCount,
	}

	if d.CPU.Total == 0 {
		return nil
	}

	d.UserPercent = float64(d.CPU.User) / float64(d.CPU.Total) * 100
	d.NicePercent = float64(d.CPU.Nice) / float64(d.CPU.Total) * 100
	d.SystemPercent = float64(d.CPU.System) / float64(d.CPU.Total) * 100
	d.IdlePercent = float64(d.CPU.Idle) / float64(d.CPU.Total) * 100

	return d
}
