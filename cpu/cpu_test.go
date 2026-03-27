package cpu

import (
	"math"
	"strings"
	"testing"
)

// Simulated /proc/stat content with 4 CPUs
const procStatContent = `cpu  10132153 290696 3084719 46828483 16683 0 25195 0 0 0
cpu0 2526892 72724 771045 11714498 4161 0 6292 0 0 0
cpu1 2525528 72571 771165 11714619 4164 0 6283 0 0 0
cpu2 2540280 72594 771137 11699410 4178 0 6309 0 0 0
cpu3 2539453 72807 771372 11699956 4180 0 6311 0 0 0
intr 330660234 0 0 0 0 0 0 0 0
`

func TestParseLine(t *testing.T) {
	fields := []string{"1000", "200", "300", "5000", "100", "0", "50", "0", "80", "20"}
	cs, err := parseLine(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.User != 1000 {
		t.Errorf("User = %d, want 1000", cs.User)
	}
	if cs.Idle != 5000 {
		t.Errorf("Idle = %d, want 5000", cs.Idle)
	}
	if cs.Guest != 80 {
		t.Errorf("Guest = %d, want 80", cs.Guest)
	}
	if cs.GuestNice != 20 {
		t.Errorf("GuestNice = %d, want 20", cs.GuestNice)
	}
	rawSum := uint64(1000 + 200 + 300 + 5000 + 100 + 0 + 50 + 0 + 80 + 20)
	expectedTotal := rawSum - 80 - 20
	if cs.Total != expectedTotal {
		t.Errorf("Total = %d, want %d", cs.Total, expectedTotal)
	}
}

func TestParseLineFewerFields(t *testing.T) {
	fields := []string{"1000", "200", "300", "5000"}
	cs, err := parseLine(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.User != 1000 {
		t.Errorf("User = %d, want 1000", cs.User)
	}
	if cs.Iowait != 0 {
		t.Errorf("Iowait = %d, want 0", cs.Iowait)
	}
}

func TestParseLineInvalidValue(t *testing.T) {
	fields := []string{"1000", "abc", "300"}
	_, err := parseLine(fields)
	if err == nil {
		t.Error("expected error for invalid value, got nil")
	}
}

func TestParseLineEmpty(t *testing.T) {
	cs, err := parseLine([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.Total != 0 {
		t.Errorf("Total = %d, want 0", cs.Total)
	}
}

func TestCollectCPUStats(t *testing.T) {
	reader := strings.NewReader(procStatContent)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.CPU.User != 10132153 {
		t.Errorf("CPU.User = %d, want 10132153", stats.CPU.User)
	}
	if stats.CPU.Nice != 290696 {
		t.Errorf("CPU.Nice = %d, want 290696", stats.CPU.Nice)
	}
	if stats.CPU.System != 3084719 {
		t.Errorf("CPU.System = %d, want 3084719", stats.CPU.System)
	}
	if stats.CPU.Idle != 46828483 {
		t.Errorf("CPU.Idle = %d, want 46828483", stats.CPU.Idle)
	}

	expectedTotal := uint64(10132153 + 290696 + 3084719 + 46828483 + 16683 + 0 + 25195 + 0 + 0 + 0)
	if stats.CPU.Total != expectedTotal {
		t.Errorf("CPU.Total = %d, want %d", stats.CPU.Total, expectedTotal)
	}

	if stats.StatCount != 10 {
		t.Errorf("StatCount = %d, want 10", stats.StatCount)
	}
	if stats.CPUCount != 4 {
		t.Errorf("CPUCount = %d, want 4", stats.CPUCount)
	}

	u := stats.CPU.Usage()
	wantUserPct := float64(10132153) / float64(expectedTotal) * 100
	if math.Abs(u.UserPercent-wantUserPct) > 0.001 {
		t.Errorf("UserPercent = %f, want %f", u.UserPercent, wantUserPct)
	}
	wantIdlePct := float64(46828483) / float64(expectedTotal) * 100
	if math.Abs(u.IdlePercent-wantIdlePct) > 0.001 {
		t.Errorf("IdlePercent = %f, want %f", u.IdlePercent, wantIdlePct)
	}
}

func TestCollectCPUStatsWithGuestTime(t *testing.T) {
	input := "cpu  1000 200 300 5000 100 0 50 0 80 20\n"
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rawSum := uint64(1000 + 200 + 300 + 5000 + 100 + 0 + 50 + 0 + 80 + 20)
	expectedTotal := rawSum - 80 - 20
	if stats.CPU.Total != expectedTotal {
		t.Errorf("CPU.Total = %d, want %d", stats.CPU.Total, expectedTotal)
	}
	if stats.CPU.Guest != 80 {
		t.Errorf("CPU.Guest = %d, want 80", stats.CPU.Guest)
	}
	if stats.CPU.GuestNice != 20 {
		t.Errorf("CPU.GuestNice = %d, want 20", stats.CPU.GuestNice)
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
	input := "cpu  1000 200 300 5000\n"
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.CPU.User != 1000 {
		t.Errorf("CPU.User = %d, want 1000", stats.CPU.User)
	}
	if stats.CPU.Idle != 5000 {
		t.Errorf("CPU.Idle = %d, want 5000", stats.CPU.Idle)
	}
	if stats.StatCount != 4 {
		t.Errorf("StatCount = %d, want 4", stats.StatCount)
	}
}

func TestCollectCPUStatsPerCore(t *testing.T) {
	input := `cpu  2000 400 600 10000 200 0 100 0 0 0
cpu0 1000 200 300 5000 100 0 50 0 0 0
cpu1 1000 200 300 5000 100 0 50 0 0 0
intr 0
`
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.CPUCount != 2 {
		t.Errorf("CPUCount = %d, want 2", stats.CPUCount)
	}
	if len(stats.Cores) != 2 {
		t.Fatalf("len(Cores) = %d, want 2", len(stats.Cores))
	}
	if stats.Cores[0].User != 1000 {
		t.Errorf("Cores[0].User = %d, want 1000", stats.Cores[0].User)
	}
	if stats.Cores[1].System != 300 {
		t.Errorf("Cores[1].System = %d, want 300", stats.Cores[1].System)
	}
	expectedCoreTotal := uint64(1000 + 200 + 300 + 5000 + 100 + 0 + 50 + 0 + 0 + 0)
	if stats.Cores[0].Total != expectedCoreTotal {
		t.Errorf("Cores[0].Total = %d, want %d", stats.Cores[0].Total, expectedCoreTotal)
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

func TestDelta(t *testing.T) {
	prev := &Stats{
		CPU: CoreStats{
			User: 1000, Nice: 200, System: 300, Idle: 5000,
			Iowait: 100, Irq: 0, Softirq: 50, Steal: 0, Guest: 0, GuestNice: 0,
			Total: 6650,
		},
		CPUCount: 4, StatCount: 10,
	}
	next := &Stats{
		CPU: CoreStats{
			User: 1200, Nice: 210, System: 350, Idle: 5400,
			Iowait: 110, Irq: 0, Softirq: 60, Steal: 0, Guest: 0, GuestNice: 0,
			Total: 7330,
		},
		CPUCount: 4, StatCount: 10,
	}

	d := Delta(prev, next)
	if d == nil {
		t.Fatal("Delta returned nil, expected non-nil")
	}

	totalDelta := uint64(7330 - 6650)
	if d.CPU.Total != totalDelta {
		t.Errorf("CPU.Total = %d, want %d", d.CPU.Total, totalDelta)
	}
	if d.CPU.User != 200 {
		t.Errorf("CPU.User = %d, want 200", d.CPU.User)
	}

	u := d.CPU.Usage()
	wantUserPct := float64(200) / float64(680) * 100
	if math.Abs(u.UserPercent-wantUserPct) > 0.001 {
		t.Errorf("UserPercent = %f, want %f", u.UserPercent, wantUserPct)
	}
	wantIdlePct := float64(400) / float64(680) * 100
	if math.Abs(u.IdlePercent-wantIdlePct) > 0.001 {
		t.Errorf("IdlePercent = %f, want %f", u.IdlePercent, wantIdlePct)
	}

	if d.CPUCount != 4 {
		t.Errorf("CPUCount = %d, want 4", d.CPUCount)
	}
}

func TestDeltaPerCore(t *testing.T) {
	core0 := CoreStats{User: 500, Idle: 2500, Total: 3000}
	core1 := CoreStats{User: 500, Idle: 2500, Total: 3000}
	prev := &Stats{
		CPU:   CoreStats{User: 1000, Idle: 5000, Total: 6000},
		Cores: []CoreStats{core0, core1},
	}

	core0next := CoreStats{User: 600, Idle: 2600, Total: 3200}
	core1next := CoreStats{User: 700, Idle: 2500, Total: 3200}
	next := &Stats{
		CPU:   CoreStats{User: 1300, Idle: 5100, Total: 6400},
		Cores: []CoreStats{core0next, core1next},
	}

	d := Delta(prev, next)
	if d == nil {
		t.Fatal("Delta returned nil")
	}
	if len(d.Cores) != 2 {
		t.Fatalf("len(Cores) = %d, want 2", len(d.Cores))
	}
	if d.Cores[0].User != 100 {
		t.Errorf("Cores[0].User = %d, want 100", d.Cores[0].User)
	}
	if d.Cores[1].User != 200 {
		t.Errorf("Cores[1].User = %d, want 200", d.Cores[1].User)
	}
	if d.Cores[0].Total != 200 {
		t.Errorf("Cores[0].Total = %d, want 200", d.Cores[0].Total)
	}
}

func TestUsage(t *testing.T) {
	cs := CoreStats{User: 200, Nice: 10, System: 50, Idle: 400, Iowait: 30, Steal: 10, Total: 700}
	u := cs.Usage()

	if math.Abs(u.UserPercent-float64(200)/700*100) > 0.001 {
		t.Errorf("UserPercent = %f", u.UserPercent)
	}
	if math.Abs(u.IowaitPercent-float64(30)/700*100) > 0.001 {
		t.Errorf("IowaitPercent = %f", u.IowaitPercent)
	}
	if math.Abs(u.StealPercent-float64(10)/700*100) > 0.001 {
		t.Errorf("StealPercent = %f", u.StealPercent)
	}
}

func TestUsageZeroTotal(t *testing.T) {
	cs := CoreStats{}
	u := cs.Usage()
	if u.UserPercent != 0 || u.IdlePercent != 0 {
		t.Errorf("expected zero Usage, got %+v", u)
	}
}

func TestUsageAllIdle(t *testing.T) {
	cs := CoreStats{Idle: 1000, Total: 1000}
	u := cs.Usage()
	if math.Abs(u.IdlePercent-100) > 0.001 {
		t.Errorf("IdlePercent = %f, want 100", u.IdlePercent)
	}
	if u.UserPercent != 0 {
		t.Errorf("UserPercent = %f, want 0", u.UserPercent)
	}
}

func TestDeltaZeroTotal(t *testing.T) {
	s := &Stats{
		CPU: CoreStats{
			User: 1000, Nice: 200, System: 300, Idle: 5000,
			Total: 6500,
		},
		CPUCount: 2,
	}
	d := Delta(s, s)
	if d != nil {
		t.Errorf("Delta with identical stats should return nil, got %+v", d)
	}
}

func TestDeltaCoreCountMismatch(t *testing.T) {
	prev := &Stats{
		CPU:   CoreStats{User: 1000, Total: 5000},
		Cores: []CoreStats{{User: 500, Total: 2500}, {User: 500, Total: 2500}},
	}
	next := &Stats{
		CPU:   CoreStats{User: 1200, Total: 5200},
		Cores: []CoreStats{{User: 600, Total: 2600}},
	}
	d := Delta(prev, next)
	if d == nil {
		t.Fatal("Delta returned nil")
	}
	if len(d.Cores) != 1 {
		t.Errorf("len(Cores) = %d, want 1 (min of prev=2, next=1)", len(d.Cores))
	}
}

func TestParseLineExtraFields(t *testing.T) {
	// Future kernels may add fields beyond the 10 we know
	fields := []string{"100", "200", "300", "400", "500", "600", "700", "800", "900", "1000", "9999", "8888"}
	cs, err := parseLine(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.GuestNice != 1000 {
		t.Errorf("GuestNice = %d, want 1000", cs.GuestNice)
	}
	// Extra fields should be ignored
	expectedTotal := uint64(100+200+300+400+500+600+700+800+900+1000) - 900 - 1000
	if cs.Total != expectedTotal {
		t.Errorf("Total = %d, want %d", cs.Total, expectedTotal)
	}
}

func TestParseLineUint64Overflow(t *testing.T) {
	// Value exceeding uint64 max
	fields := []string{"99999999999999999999999"}
	_, err := parseLine(fields)
	if err == nil {
		t.Error("expected error for uint64 overflow, got nil")
	}
}

func TestCollectCPUStatsInvalidCoreLine(t *testing.T) {
	input := "cpu  1000 200 300 5000 100 0 50 0 0 0\ncpu0 abc 200 300 5000 100 0 50 0 0 0\n"
	reader := strings.NewReader(input)
	_, err := collectCPUStats(reader)
	if err == nil {
		t.Error("expected error for invalid core line, got nil")
	}
}

func TestCollectCPUStatsPerCoreGuestSubtraction(t *testing.T) {
	input := `cpu  1000 200 300 5000 100 0 50 0 80 20
cpu0 500 100 150 2500 50 0 25 0 40 10
`
	reader := strings.NewReader(input)
	stats, err := collectCPUStats(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Cores[0].Guest != 40 {
		t.Errorf("Cores[0].Guest = %d, want 40", stats.Cores[0].Guest)
	}
	if stats.Cores[0].GuestNice != 10 {
		t.Errorf("Cores[0].GuestNice = %d, want 10", stats.Cores[0].GuestNice)
	}
	rawSum := uint64(500 + 100 + 150 + 2500 + 50 + 0 + 25 + 0 + 40 + 10)
	expectedTotal := rawSum - 40 - 10
	if stats.Cores[0].Total != expectedTotal {
		t.Errorf("Cores[0].Total = %d, want %d", stats.Cores[0].Total, expectedTotal)
	}
}

func TestUsagePerCoreDelta(t *testing.T) {
	prev := &Stats{
		CPU:   CoreStats{User: 2000, System: 600, Idle: 10000, Total: 12600},
		Cores: []CoreStats{{User: 1000, System: 300, Idle: 5000, Total: 6300}},
	}
	next := &Stats{
		CPU:   CoreStats{User: 2200, System: 650, Idle: 10500, Total: 13350},
		Cores: []CoreStats{{User: 1100, System: 325, Idle: 5250, Total: 6675}},
	}
	d := Delta(prev, next)
	if d == nil {
		t.Fatal("Delta returned nil")
	}
	u := d.Cores[0].Usage()
	// Core delta: User=100, System=25, Idle=250, Total=375
	wantUserPct := float64(100) / float64(375) * 100
	if math.Abs(u.UserPercent-wantUserPct) > 0.001 {
		t.Errorf("Cores[0] UserPercent = %f, want %f", u.UserPercent, wantUserPct)
	}
}
