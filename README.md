# cpustat

[![CI](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml/badge.svg)](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hidetzu/cpustat.svg)](https://pkg.go.dev/github.com/hidetzu/cpustat)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidetzu/cpustat)](https://goreportcard.com/report/github.com/hidetzu/cpustat)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A small Go library to get Linux CPU statistics from `/proc/stat`.

## Install

```sh
go install github.com/hidetzu/cpustat@latest
```

## Supported OS

Linux Only

## Examples

### Snapshot

```go
package main

import (
	"fmt"
	"log"

	"github.com/hidetzu/cpustat/cpu"
)

func main() {
	stats, err := cpu.Get()
	if err != nil {
		log.Fatal(err)
	}

	u := stats.CPU.Usage()
	fmt.Printf("user%%\tnice%%\tsystem%%\tidle%%\n")
	fmt.Printf("%.1f\t%.1f\t%.1f\t%.1f\n",
		u.UserPercent, u.NicePercent, u.SystemPercent, u.IdlePercent,
	)
}
```

### Delta-based usage

```go
prev, _ := cpu.Get()
time.Sleep(5 * time.Second)
next, _ := cpu.Get()

if d := cpu.Delta(prev, next); d != nil {
	u := d.CPU.Usage()
	fmt.Printf("user: %.1f%%  idle: %.1f%%\n", u.UserPercent, u.IdlePercent)
}
```

### Per-core stats

```go
stats, _ := cpu.Get()
for i, core := range stats.Cores {
	u := core.Usage()
	fmt.Printf("cpu%d: user=%.1f%% idle=%.1f%%\n", i, u.UserPercent, u.IdlePercent)
}
```

## License

MIT
