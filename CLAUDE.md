# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Linux-only Go library that reads CPU statistics from `/proc/stat`. Designed to be a small, focused library — not a CLI tool, TUI, or cross-platform system monitor. The `main.go` is a minimal demo only.

## Build Commands

```sh
# Build via Docker (outputs to bin/cpustat)
make

# Build locally
go build -o bin/cpustat .

# Run tests
go test ./...

# Run tests with race detection
go test ./... -v -race
```

## Architecture

- **`cpu/cpu.go`** — Core library with three key types:
  - `CoreStats` — Raw tick counters for one CPU line (aggregate or per-core). Guest/guest_nice are subtracted from total per Linux kernel convention.
  - `Stats` — Snapshot of `/proc/stat` containing `CPU` (aggregate), `Cores` (per-core `[]CoreStats`), `CPUCount`, and `StatCount`.
  - `Usage` — Percentage breakdown (user, nice, system, idle, iowait, steal) computed via `CoreStats.Usage()`.
  - Key functions: `Get()` returns a snapshot, `Delta(prev, next)` computes tick deltas between two snapshots.
  - Internal: `parseLine()` parses a single `/proc/stat` CPU line, reused for both aggregate and per-core lines.
- **`main.go`** — Demo CLI. Polls `cpu.Get()` on a 5-second ticker, prints delta-based usage via `CPU.Usage()`. Graceful shutdown on SIGTERM/SIGINT.
- **`build/Dockerfile`** — Multi-stage build container using `golang:1.20`.

## Platform Constraint

Linux only — depends on `/proc/stat`. Will not compile or run correctly on macOS/Windows.
