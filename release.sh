#!/bin/bash

# release.sh - Automates the release process for Timan (Formula + Cask)
# Usage: ./release.sh v1.0.1

set -e

if [ -z "$1" ]; then
    echo "Usage: ./release.sh <version> (e.g., v1.0.1)"
    exit 1
fi

VERSION=$1
VERSION_NUM=${VERSION#v}
APP_NAME="Timan"
BUNDLE_ID="com.timan.timer"
GOPATH_BIN=$(go env GOPATH)/bin

echo "🚀 Preparing Professional Release for $VERSION..."

# 1. Ensure working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Error: Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# 2. Build and Package the .app bundle
echo "📦 Building native macOS .app bundle..."
rm -rf "$APP_NAME.app" "$APP_NAME.zip"
$GOPATH_BIN/fyne package -os darwin -icon icon.png -name "$APP_NAME" -id "$BUNDLE_ID"

# 3. Create a Zip for the Cask
echo "🤐 Zipping $APP_NAME.app for distribution..."
zip -r "$APP_NAME.zip" "$APP_NAME.app"
CASK_SHA=$(shasum -a 256 "$APP_NAME.zip" | awk '{print $1}')

# 4. Update the CLI Formula
echo "📝 Updating Homebrew Formula..."
sed -i '' "s/version \".*\"/version \"$VERSION_NUM\"/g" Formula/timan.rb
sed -i '' "s/tags\/v.*\.tar\.gz/tags\/$VERSION.tar.gz/g" Formula/timan.rb

# 5. Push code and tags to GitHub
git add .
git commit -m "chore: release $VERSION" || true
git push origin master
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

# 6. Upload Assets to GitHub Release
echo "📤 Creating GitHub Release and uploading $APP_NAME.zip..."
gh release create "$VERSION" "$APP_NAME.zip" --title "$VERSION" --notes "Timan $VERSION HUD Release"

# 7. Calculate Formula SHA (from the auto-generated GitHub source tarball)
echo "⏳ Waiting for source archive..."
sleep 5
ARCHIVE_URL="https://github.com/kena421/timan/archive/refs/tags/$VERSION.tar.gz"
curl -L "$ARCHIVE_URL" -o "source.tar.gz"
SOURCE_SHA=$(shasum -a 256 source.tar.gz | awk '{print $1}')
rm source.tar.gz

# 8. Update Formula and Cask in main repo
sed -i '' "s/sha256 \".*\"/sha256 \"$SOURCE_SHA\"/g" Formula/timan.rb
git add Formula/timan.rb
git commit --amend --no-edit
git push origin master --force
git tag -f "$VERSION"
git push origin "$VERSION" --force

# 9. Update the Homebrew Tap (Formula AND Cask)
echo "🍺 Updating Homebrew Tap (kena421/homebrew-tap)..."
TAP_DIR="/tmp/timan-tap-release"
rm -rf "$TAP_DIR"
gh repo clone kena421/homebrew-tap "$TAP_DIR"

# Update Formula
mkdir -p "$TAP_DIR/Formula"
cp Formula/timan.rb "$TAP_DIR/Formula/timan.rb"

# Update/Create Cask
mkdir -p "$TAP_DIR/Casks"
cat <<EOF > "$TAP_DIR/Casks/timan.rb"
cask "timan" do
  version "$VERSION_NUM"
  sha256 "$CASK_SHA"

  url "https://github.com/kena421/timan/releases/download/v#{version}/Timan.zip"
  name "Timan"
  desc "Minimalist HUD Time Manager"
  homepage "https://github.com/kena421/timan"

  app "Timan.app"

  zap trash: "~/.config/timan"
end
EOF

cd "$TAP_DIR"
git add .
git commit -m "feat: release $VERSION (Formula + Cask)"
git push origin main
rm -rf "$TAP_DIR"

# Cleanup local build
rm -rf "$APP_NAME.app" "$APP_NAME.zip"

echo "🎉 Professional Release $VERSION complete!"
echo "Users can now install the App via: brew install --cask timan"
