# Contributing to Timan

Thank you for your interest in contributing to **Timan**! We welcome improvements to the core engine, UI enhancements, and better platform support.

## How to Contribute

1. **Fork the repository** on GitHub.
2. **Clone your fork** to your local machine.
3. **Create a new branch** for your feature or bugfix:
   ```bash
   git checkout -b feature/amazing-new-feature
   ```
4. **Make your changes**. Ensure your code follows standard Go formatting (`go fmt`).
5. **Commit your changes**:
   ```bash
   git commit -m "feat: add amazing new feature"
   ```
6. **Push to your branch**:
   ```bash
   git push origin feature/amazing-new-feature
   ```
7. **Open a Pull Request** against the main repository.

## Development Setup

Timan requires Go 1.21+ and a macOS environment for the native HUD features.

### Building Locally

To build and run the app during development without installing:
```bash
go run main.go
```

### UI Development

The UI is built using [Fyne](https://fyne.io). If you are adding new widgets or changing layouts, please test across different window sizes.

## Coding Standards

- **Go Doc**: Please add documentation comments to all new exported types and functions.
- **Consistency**: Follow the existing project structure (Internal/Domain, Internal/Engine, etc.).

## Questions?

Feel free to open an issue for any questions or architectural discussions!
