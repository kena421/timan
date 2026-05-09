# Interview Timer (Go + Fyne)

A minimalist, "Always-on-Top", frameless interview timer for macOS.

## New Features (v2)
- **Click to Start/Pause**: The timer starts in a paused state (grey). Click the window to start/pause.
- **Phase-Based Timing**: Automatically transitions through customizable interview phases.
- **Progress Tracking**: Mini progress bar shows progress within the current phase.
- **Dynamic Color**: Changes from Grey (Paused) to Green (Active) to Orange (Last minute).
- **Custom Configuration**: Right-click to define your own phases (e.g., `Intro:5,Coding:30,Review:5`).

## Setup & Quick Installation
Ensure you have [Go](https://go.dev/dl/) installed. Then run the automated install script:

```bash
chmod +x install.sh
./install.sh
```

This script will:
1. Build the optimized binary.
2. Package it as a native macOS `.app` bundle with a custom icon.
3. Install it to your `/Applications` folder.
4. Install the CLI tool to `/usr/local/bin/timer`.

## Usage
- **Launch**: Open **Timer** from your Applications folder or run `timer` in terminal.
- **Move**: Click and drag the timer anywhere on the screen.
- **Switch Mode**: Right-click and select **Toggle Mode (Up/Down)**.
- **Reset**: Right-click and select **Reset**.
- **Quit**: Right-click and select **Quit**.

---

*Note: This application uses native macOS APIs (Cocoa) via CGO and is specifically optimized for macOS.*
