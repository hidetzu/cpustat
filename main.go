package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hidetzu/cpustat/cpu"
)

func cancelContextWithSigterm(ctx context.Context) context.Context {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-signals
		cancel()
	}()
	return ctx
}

func main() {
	ctx := cancelContextWithSigterm(context.Background())
	startWork(ctx)
}

// startWork polls CPU stats every 5 seconds and prints delta-based usage.
func startWork(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var prev *cpu.Stats
	for {
		cur, err := cpu.Get()
		if err != nil {
			fmt.Printf("failed to get cpu stats: %s\n", err)
		} else if prev != nil {
			if d := cpu.Delta(prev, cur); d != nil {
				fmt.Printf("user%%\tnice%%\tsystem%%\tidle%%\n")
				fmt.Printf("%.1f\t%.1f\t%.1f\t%.1f\n",
					d.UserPercent,
					d.NicePercent,
					d.SystemPercent,
					d.IdlePercent,
				)
			}
		}
		prev = cur

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}
