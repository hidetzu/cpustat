# cpustat

[![CI](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml/badge.svg)](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hidetzu/cpustat.svg)](https://pkg.go.dev/github.com/hidetzu/cpustat)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidetzu/cpustat)](https://goreportcard.com/report/github.com/hidetzu/cpustat)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A small Go library to get Linux CPU statistics from `/proc/stat`.

`cpu.Get()` returns a snapshot of raw tick counters. `cpu.Delta(prev, next)` computes the difference between two snapshots, and `CoreStats.Usage()` converts tick deltas into percentages. Both aggregate and per-core stats are supported.

The `main.go` in this repository is a minimal demo; the library itself is the primary deliverable.

## Install

```sh
go get github.com/hidetzu/cpustat
```

## Supported OS

Linux Only

## Examples

### Snapshot

```go
stats, err := cpu.Get()
if err != nil {
	log.Fatal(err)
}
u := stats.CPU.Usage()
fmt.Printf("user=%.1f%% system=%.1f%% idle=%.1f%%\n",
	u.UserPercent, u.SystemPercent, u.IdlePercent)
```

### Delta-based usage

```go
prev, err := cpu.Get()
if err != nil {
	log.Fatal(err)
}
time.Sleep(5 * time.Second)
next, err := cpu.Get()
if err != nil {
	log.Fatal(err)
}

if d := cpu.Delta(prev, next); d != nil {
	u := d.CPU.Usage()
	fmt.Printf("user=%.1f%% system=%.1f%% idle=%.1f%%\n",
		u.UserPercent, u.SystemPercent, u.IdlePercent)
}
```

### Per-core stats

```go
stats, err := cpu.Get()
if err != nil {
	log.Fatal(err)
}
for i, core := range stats.Cores {
	u := core.Usage()
	fmt.Printf("cpu%d: user=%.1f%% idle=%.1f%%\n",
		i, u.UserPercent, u.IdlePercent)
}
```

## Testing

Tests cover the parser (`parseLine`), snapshot collection, delta calculation, `Usage()` computation, and edge cases (fewer fields, extra fields, uint64 overflow, invalid values, core count mismatch, guest time subtraction).

```sh
go test ./... -v -race
```

## License

MIT
