# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Linux-only Go library and CLI tool that reads CPU statistics from `/proc/stat`. The `cpu` package exposes a `Get()` function returning a `Stats` struct with raw counters and percentage breakdowns. The `main` package runs a polling loop (every 5 seconds) that prints CPU usage until interrupted by SIGTERM/SIGINT.

## Build Commands

```sh
# Build via Docker (outputs to bin/cpustat)
make

# Build locally
go build -o bin/cpustat .

# Run tests
go test ./...
```

## Architecture

- **`cpu/cpu.go`** — Core library. Parses `/proc/stat` to populate `Stats` (user, nice, system, idle, iowait, irq, softirq, steal, guest, guest_nice, total, per-CPU count, percentages). Guest/guest_nice are subtracted from total to avoid double-counting per Linux kernel convention.
- **`main.go`** — CLI entry point. Polls `cpu.Get()` on a 5-second ticker with graceful shutdown via context cancellation on SIGTERM/SIGINT.
- **`build/Dockerfile`** — CentOS Stream 9 based build container.

## Platform Constraint

Linux only — depends on `/proc/stat`. Will not compile or run correctly on macOS/Windows.
