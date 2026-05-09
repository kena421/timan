# Interview Timer (Go + Fyne)

A minimalist, "Always-on-Top", frameless interview timer for macOS.

## Features
- **Frameless Window**: No title bar or borders.
- **Always on Top**: Stays visible above full-screen IDEs and browsers.
- **Dual Mode**: Count Up (default) or Count Down (from 45:00).
- **Draggable**: Click and drag anywhere on the timer to reposition.
- **Context Menu**: Right-click for Toggle Mode, Reset, and Quit.
- **Visual Cues**: Green for active, Red when exceeding limit/reaching zero.
- **Lightweight**: Optimized binary size using stripped debug symbols.

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
