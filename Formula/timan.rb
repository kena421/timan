class Timan < Formula
  desc "Minimalist HUD Time Manager for macOS"
  homepage "https://github.com/kena421/timan"
  url "https://github.com/kena421/timan/archive/refs/tags/v1.0.1.tar.gz"
  version "1.0.1"
  sha256 "bae69979f7dc3859fed1a0bd1280ef9af6fa1f8b5db5efc10637097956923d2d"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build the optimized CLI binary using standard Homebrew Go arguments
    system "go", "build", *std_go_args(ldflags: "-s -w"), "main.go"
    
    # Note: Homebrew formulae are primarily for CLI tools.
    # To install the .app bundle, a Homebrew Cask is typically used.
    # This formula installs the 'timan' command which can be used to launch the HUD.
  end

  test do
    # Simple version check
    system "#{bin}/timan", "--help"
  end
end
