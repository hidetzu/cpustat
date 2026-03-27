package cpu

import (
	"math"
	"strings"
	"testing"
)

// Simulated /proc/stat content with 4 CPUs
const procStatContent = `cpu  10132153 290696 3084719 46828483 16683 0 25195 0 0 0
cpu0 2526892 72724 771045 11714498 4161 0 6292 0 0 0
cpu1 2525528 72## 771165 11714619 4164 0 6283 0 0 0
cpu2 2540280 72594 771137 11699410 4178 0 6309 0 0 0
cpu3 2539453 72807 771372 11699956 4180 0 6311 0 0 0
intr 330660234 0 0 0 0 0 0 0 0
`

func TestCollectCPUStats(t *testing.T) {
	reader := strings.NewReader(procStatContent)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify raw values from the first "cpu" line
	if stats.User != 10132153 {
		t.Errorf("User = %d, want 10132153", stats.User)
	}
	if stats.Nice != 290696 {
		t.Errorf("Nice = %d, want 290696", stats.Nice)
	}
	if stats.System != 3084719 {
		t.Errorf("System = %d, want 3084719", stats.System)
	}
	if stats.Idle != 46828483 {
		t.Errorf("Idle = %d, want 46828483", stats.Idle)
	}
	if stats.Iowait != 16683 {
		t.Errorf("Iowait = %d, want 16683", stats.Iowait)
	}
	if stats.Irq != 0 {
		t.Errorf("Irq = %d, want 0", stats.Irq)
	}
	if stats.Softirq != 25195 {
		t.Errorf("Softirq = %d, want 25195", stats.Softirq)
	}
	if stats.Steal != 0 {
		t.Errorf("Steal = %d, want 0", stats.Steal)
	}
	if stats.Guest != 0 {
		t.Errorf("Guest = %d, want 0", stats.Guest)
	}
	if stats.GuestNice != 0 {
		t.Errorf("GuestNice = %d, want 0", stats.GuestNice)
	}

	// Total = sum of all fields - Guest - GuestNice
	expectedTotal := uint64(10132153 + 290696 + 3084719 + 46828483 + 16683 + 0 + 25195 + 0 + 0 + 0)
	if stats.Total != expectedTotal {
		t.Errorf("Total = %d, want %d", stats.Total, expectedTotal)
	}

	if stats.StatCount != 10 {
		t.Errorf("StatCount = %d, want 10", stats.StatCount)
	}

	if stats.CPUCount != 4 {
		t.Errorf("CPUCount = %d, want 4", stats.CPUCount)
	}

	// Verify percentages
	wantUserPct := float64(10132153) / float64(expectedTotal) * 100
	if math.Abs(stats.UserPercent-wantUserPct) > 0.001 {
		t.Errorf("UserPercent = %f, want %f", stats.UserPercent, wantUserPct)
	}
	wantIdlePct := float64(46828483) / float64(expectedTotal) * 100
	if math.Abs(stats.IdlePercent-wantIdlePct) > 0.001 {
		t.Errorf("IdlePercent = %f, want %f", stats.IdlePercent, wantIdlePct)
	}
}

func TestCollectCPUStatsWithGuestTime(t *testing.T) {
	// Guest and GuestNice are non-zero; they should be subtracted from Total
	input := "cpu  1000 200 300 5000 100 0 50 0 80 20\n"
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rawSum := uint64(1000 + 200 + 300 + 5000 + 100 + 0 + 50 + 0 + 80 + 20)
	expectedTotal := rawSum - 80 - 20 // subtract Guest and GuestNice
	if stats.Total != expectedTotal {
		t.Errorf("Total = %d, want %d (raw sum %d minus Guest 80 and GuestNice 20)", stats.Total, expectedTotal, rawSum)
	}

	if stats.Guest != 80 {
		t.Errorf("Guest = %d, want 80", stats.Guest)
	}
	if stats.GuestNice != 20 {
		t.Errorf("GuestNice = %d, want 20", stats.GuestNice)
	}
}

func TestCollectCPUStatsEmptyInput(t *testing.T) {
	reader := strings.NewReader("")
	_, err := collectCPUStats(reader)
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestCollectCPUStatsInvalidValue(t *testing.T) {
	input := "cpu  1000 abc 300 5000 100 0 50 0 0 0\n"
	reader := strings.NewReader(input)
	_, err := collectCPUStats(reader)
	if err == nil {
		t.Error("expected error for invalid value, got nil")
	}
}

func TestCollectCPUStatsFewerFields(t *testing.T) {
	// Old kernels may have fewer fields (e.g., 4 instead of 10)
	input := "cpu  1000 200 300 5000\n"
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.User != 1000 {
		t.Errorf("User = %d, want 1000", stats.User)
	}
	if stats.Idle != 5000 {
		t.Errorf("Idle = %d, want 5000", stats.Idle)
	}
	if stats.Iowait != 0 {
		t.Errorf("Iowait = %d, want 0 (not present in input)", stats.Iowait)
	}
	if stats.StatCount != 4 {
		t.Errorf("StatCount = %d, want 4", stats.StatCount)
	}
}

func TestCollectCPUStatsNoCPULines(t *testing.T) {
	input := "cpu  1000 200 300 5000 100 0 50 0 0 0\nintr 0\n"
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.CPUCount != 0 {
		t.Errorf("CPUCount = %d, want 0", stats.CPUCount)
	}
}
