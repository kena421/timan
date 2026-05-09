#!/bin/bash

# Exit on error
set -e

APP_NAME="Timer"
BINARY_NAME="timer"
GOPATH_BIN=$(go env GOPATH)/bin

echo "🚀 Starting installation for $APP_NAME..."

# 1. Check for Fyne tool
if [ ! -f "$GOPATH_BIN/fyne" ]; then
    echo "📦 Installing Fyne packaging tool..."
    go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# 2. Build the binary
echo "🛠 Building optimized binary..."
go build -ldflags="-s -w" -o $BINARY_NAME main.go

# 3. Handle icon (Convert to PNG if needed)
if [ -f "icon.png" ]; then
    echo "🎨 Preparing icon..."
    # Ensure it's a real PNG (sips is built-in on Mac)
    sips -s format png icon.png --out icon_processed.png > /dev/null 2>&1
    ICON_FLAG="-icon icon_processed.png"
else
    ICON_FLAG=""
fi

# 4. Package as .app bundle
echo "📦 Packaging as .app bundle..."
rm -rf "$APP_NAME.app" "InterviewTimer.app"
$GOPATH_BIN/fyne package -os darwin $ICON_FLAG -name $APP_NAME -id com.interview.timer

# 5. Build the CLI binary (after packaging to avoid deletion)
echo "🛠 Building optimized CLI binary..."
go build -ldflags="-s -w" -o $BINARY_NAME main.go

# 6. Install to /Applications
echo "🚚 Installing to /Applications..."
sudo rm -rf "/Applications/$APP_NAME.app"
sudo mv "$APP_NAME.app" "/Applications/"

# 7. Install CLI tool (Optional)
echo "🔗 Installing CLI tool to /usr/local/bin..."
sudo mv $BINARY_NAME /usr/local/bin/$BINARY_NAME

# Cleanup
rm -f icon_processed.png

echo "✅ Installation complete!"
echo "✨ You can now find '$APP_NAME' in your Applications folder or run '$BINARY_NAME' in your terminal."
