# cpustat

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
