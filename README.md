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

## Setup
Ensure you have [Go](https://go.dev/dl/) installed.

```bash
# Initialize and fetch dependencies
go mod tidy
```

## Build
To create a small, optimized executable:

```bash
go build -ldflags="-s -w" -o timer main.go
```

## Usage
- **Run**: `./timer`
- **Move**: Click and drag the timer to your preferred corner.
- **Switch Mode**: Right-click and select "Toggle Mode".
- **Reset**: Right-click and select "Reset".
- **Quit**: Right-click and select "Quit".

---

*Note: This application uses native macOS APIs (Cocoa) via CGO to achieve the frameless and always-on-top behavior. It is specifically optimized for macOS.*
