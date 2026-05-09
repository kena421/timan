class Timan < Formula
  desc "Minimalist HUD Time Manager for macOS"
  homepage "https://github.com/kena421/timan"
  url "https://github.com/kena421/timan/archive/refs/heads/master.tar.gz"
  version "1.0.0"
  sha256 "562991532755b743d5e6ac8b561f082bda458195a75f831f55d5b29afaf253c5"
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
