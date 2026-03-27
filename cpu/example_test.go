package cpu_test

import (
	"fmt"
	"time"

	"github.com/hidetzu/cpustat/cpu"
)

func ExampleGet() {
	stats, err := cpu.Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	u := stats.CPU.Usage()
	fmt.Printf("user=%.1f%% idle=%.1f%%\n", u.UserPercent, u.IdlePercent)
}

func ExampleDelta() {
	prev, err := cpu.Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(1 * time.Second)
	next, err := cpu.Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	if d := cpu.Delta(prev, next); d != nil {
		u := d.CPU.Usage()
		fmt.Printf("user=%.1f%% idle=%.1f%%\n", u.UserPercent, u.IdlePercent)
	}
}
