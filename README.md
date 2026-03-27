# cpustat

[![CI](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml/badge.svg)](https://github.com/hidetzu/cpustat/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hidetzu/cpustat.svg)](https://pkg.go.dev/github.com/hidetzu/cpustat)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidetzu/cpustat)](https://goreportcard.com/report/github.com/hidetzu/cpustat)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a library to get linux cpu stat.

## install

```sh
go get github.com/hidetzu/cpustat
```

## Supported OS

Linux Only

## Examples

```go

package main

func main() {
	cpus, _ := cpu.Get()
	fmt.Printf("user%%\tnice%%\tsystem%%\tidle%%\n")
	fmt.Printf("%v\t%v\t%v\t%v\n",
		cpus.UserPercent,
		cpus.NicePercent,
		cpus.SystemPercent,
		cpus.IdlePercent,
	)
}
```

## License

MIT
