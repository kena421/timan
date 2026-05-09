# Timan: The Professional HUD Time Manager 🚀

**Timan** (Time Manager) is a high-performance, minimalist macOS utility designed for professionals who need precision time management without distraction. Built with a sleek **Industrial HUD** aesthetic, Timan stays out of your way while keeping you perfectly on schedule.

![Timan Interface](https://raw.githubusercontent.com/kena421/timan/master/assets/demo.png)

## ✨ Premium Features

### 🏢 Ultra-Persistent HUD
*   **Spotlight-Level Persistence**: Locked at `NSPopUpMenuWindowLevel`, the timer stays above almost every other window, including Full Screen apps and system overlays.
*   **Space-Aware**: Automatically follows you across all macOS Spaces and Desktops. It never gets lost when you switch contexts.
*   **Unconstrained Movement**: Drag the HUD anywhere, including "non-safe" areas like over the menu bar or under the notch.

### 🔒 Screen Sharing Privacy (The "Invisible" Timer)
*   **Native Privacy Mode**: With one click, make the timer window **invisible to screen sharing, recording, and screenshots**.
*   **Stealth Mode**: You see the timer, but your audience on Zoom, Google Meet, or Slack sees nothing but your background. Perfect for interviews and presentations.

### 🏗️ Structured Event Management
*   **Phase-Based Timing**: Design complex sessions with multiple phases (e.g., *Intro -> Technical -> Q&A*).
*   **Smart Persistence**: Timan remembers your last-used profile and restores it instantly on launch.
*   **Quick Timer**: Need a simple countdown? Use the "Simple Quick Timer" for one-off tasks.

### 🎨 Industrial Aesthetic
*   **Neon Feedback**: High-contrast, state-aware coloring:
    *   🟢 **Emerald Green**: Timer is active and healthy.
    *   ⚪ **Ghost Green**: Timer is paused.
    *   🔴 **Neon Red**: Alert threshold reached (Warning).
*   **Compact Footprint**: Optimized for utility-first performance with zero window decorations and a tiny screen footprint.

## 🚀 Installation

Ensure you have [Go](https://go.dev/dl/) installed. Then run the automated installation script:

```bash
git clone https://github.com/kena421/timan.git
cd timan
chmod +x install.sh
./install.sh
```

The script handles everything:
1. Builds the optimized binary.
2. Packages it as a native macOS `.app` bundle (**Timan.app**).
3. Installs it to your `/Applications` folder.
4. Installs the CLI tool to `/usr/local/bin/timan`.

## 🛠 Usage

*   **Launch**: Open **Timan** from Applications or run `timan` in terminal.
*   **Controls**: Use the flushed control box in the top-right:
    *   **Play/Pause**: Toggle the timer.
    *   **Reset**: Restart the current session.
    *   **Privacy (Eye Icon)**: Toggle visibility to screen sharing.
    *   **Settings (Gear Icon)**: Open the Event Library and Designer.
*   **Movement**: Click and drag anywhere on the HUD background to reposition it.

## 📂 Configuration

Timan stores your event templates and app state in `~/.config/timan/`. 
*   `events.json`: Your custom event templates.
*   `state.json`: Remembers your last-used profile.

## 🤝 Contributing

Contributions are welcome! Whether it's a bug fix, a new feature, or documentation improvements, feel free to open a PR.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Note: This application is specifically optimized for macOS and uses CGO to interact with native Cocoa APIs for its unique windowing behavior.*
