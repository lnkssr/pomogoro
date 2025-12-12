# Pomogoro – Clean Architecture Pomodoro Timer (CLI)

Pomogoro is a fully terminal-based Pomodoro timer implemented using Clean Architecture principles.
The application provides a modular domain/service/interface structure, a subcommand-based CLI, optional interactive confirmation mode, centered terminal UI, and acoustic notifications.

## Features

* Clean Architecture (domain, service, kernel, interface layers)
* Modern Go subcommand-style CLI (`pomogoro run ...`)
* Configurable work/break durations and cycle count
* Optional “confirmation mode” requiring keypress between phases
* Centered terminal UI with live updates
* Audio notifications via `ffplay` (optional)

## Requirements

* Go 1.21+
* (Optional) FFmpeg installed for sound playback

## Build & Run

```bash
git clone https://github.com/lnkssr/pomogoro
cd pomogoro

go mod download 
go build -o pomogoro ./cmd

./pomogoro run -w 25 -b 5 -c 4
```

## Help (more usage info)

```bash
./pomogoro help
```

