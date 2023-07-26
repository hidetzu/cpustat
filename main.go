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

// startWork performs a task every 60 seconds until the context is done.
func startWork(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		// Do work here so we don't need duplicate calls. It will run immediately, and again every minute as the loop continues.
		if err := work(ctx); err != nil {
			fmt.Printf("failed to do work: %s", err)
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func work(ctx context.Context) error {
	cpus, _ := cpu.Get()
	fmt.Printf("user%%\tnice%%\tsystem%%\tidle%%\n")
	fmt.Printf("%v\t%v\t%v\t%v\n",
		cpus.UserPercent,
		cpus.NicePercent,
		cpus.SystemPercent,
		cpus.IdlePercent,
	)
	return nil
}
