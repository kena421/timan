#!/bin/bash

# release.sh - Automates the release process for Timan
# Usage: ./release.sh v1.0.1

set -e

if [ -z "$1" ]; then
    echo "Usage: ./release.sh <version> (e.g., v1.0.1)"
    exit 1
fi

VERSION=$1
# Strip leading 'v' for internal versioning if needed
VERSION_NUM=${VERSION#v}

echo "🚀 Preparing release for $VERSION..."

# 1. Ensure working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Error: Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# 2. Update version in Homebrew Formula
echo "📝 Updating Homebrew Formula version to $VERSION..."
sed -i '' "s/version \".*\"/version \"$VERSION_NUM\"/g" Formula/timan.rb
sed -i '' "s/tags\/v.*\.tar\.gz/tags\/$VERSION.tar.gz/g" Formula/timan.rb

# 3. Temporarily commit to get the tag hash
git add Formula/timan.rb
git commit -m "chore: release $VERSION" || true
git tag -a "$VERSION" -m "Release $VERSION"

echo "📤 Pushing $VERSION to GitHub..."
git push origin master
git push origin "$VERSION"

# 4. Calculate new SHA256 from the GitHub archive
echo "⏳ Waiting for GitHub to generate the archive..."
sleep 5 # Give GitHub a moment to process the tag

ARCHIVE_URL="https://github.com/kena421/timan/archive/refs/tags/$VERSION.tar.gz"
echo "🔍 Downloading $ARCHIVE_URL to calculate checksum..."
curl -L "$ARCHIVE_URL" -o "release.tar.gz"
NEW_SHA=$(shasum -a 256 release.tar.gz | awk '{print $1}')
rm release.tar.gz

echo "✅ New SHA256: $NEW_SHA"

# 5. Update the formula with the real SHA256
sed -i '' "s/sha256 \".*\"/sha256 \"$NEW_SHA\"/g" Formula/timan.rb
git add Formula/timan.rb
git commit --amend --no-edit
git push origin master --force
git tag -f "$VERSION"
git push origin "$VERSION" --force

# 6. Update the Homebrew Tap repository
echo "🍺 Updating Homebrew Tap (kena421/homebrew-tap)..."
TAP_DIR="/tmp/timan-tap-release"
rm -rf "$TAP_DIR"
gh repo clone kena421/homebrew-tap "$TAP_DIR"
cp Formula/timan.rb "$TAP_DIR/Formula/timan.rb"

cd "$TAP_DIR"
git add .
git commit -m "feat: release $VERSION"
git push origin main
rm -rf "$TAP_DIR"

echo "🎉 Release $VERSION complete!"
echo "Users can now update via: brew update && brew upgrade timan"
