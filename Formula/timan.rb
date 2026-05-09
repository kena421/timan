class Timan < Formula
  desc "Minimalist HUD Time Manager for macOS"
  homepage "https://github.com/kena421/timan"
  url "https://github.com/kena421/timan/archive/refs/tags/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "d6fc7f1e303305542d93c9f542306c50f8906570280065963b3b08f011956aee"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build the optimized CLI binary
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"timan", "main.go"
    
    # Note: Homebrew formulae are primarily for CLI tools.
    # To install the .app bundle, a Homebrew Cask is typically used.
    # This formula installs the 'timan' command which can be used to launch the HUD.
  end

  test do
    # Simple version check
    system "#{bin}/timan", "--help"
  end
end
