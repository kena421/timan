# Timan: The Ultimate Time Manager 🚀

**Timan** (Time Manager) is a minimalist, "Always-on-Top", frameless event timer designed to help you navigate through complex, multi-phase activities with precision.

Whether you're conducting technical interviews, delivering a keynote presentation, managing a workshop, or following a strict study routine, Timan keeps you on track without cluttering your workspace.

![Timan Demo Placeholder](https://via.placeholder.com/600x400?text=Timan+HUD+Timer+Interface)

## ✨ Features

- **Phase-Based Timing**: Automatically transition through customizable event phases (e.g., Intro, Design, Q&A).
- **HUD-Style Interface**: A sleek, borderless, and semi-transparent window that stays above all others.
- **Event Library**: Save your structured "Events" as templates for repeatable workflows.
- **Dynamic Feedback**: Visual cues change color as you approach the end of a session or phase.
- **Smart Alerts**: Set custom warning thresholds (in minutes or percentages) to get notified when time is running low.
- **Optimized for macOS**: Leverages native macOS APIs for a premium, integrated feel.

## 🚀 Installation

Ensure you have [Go](https://go.dev/dl/) installed. Then run the automated installation script:

```bash
chmod +x install.sh
./install.sh
```

The script will:
1. Build the optimized binary.
2. Package it as a native macOS `.app` bundle (**Timan.app**).
3. Install it to your `/Applications` folder.
4. Install the CLI tool to `/usr/local/bin/timan`.

## 🛠 Usage

- **Launch**: Open **Timan** from your Applications folder or run `timan` in your terminal.
- **Start/Pause**: Simply click the timer window to toggle the active state.
- **Manage Events**: Right-click (or click the settings icon) to open the **Event Library**.
- **Configure**: Use the **Event Designer** to create complex sequences with specific durations or percentage-based splits.
- **Reset**: Quickly restart the current session or a specific phase from the interface.

## 📂 Configuration

Timan stores your event templates in `~/.config/timan/events.json`. You can share this file with others to sync your event structures.

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Note: This application is specifically optimized for macOS and uses CGO to interact with native Cocoa APIs.*
