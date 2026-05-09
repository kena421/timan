class Timan < Formula
  desc "Minimalist HUD Time Manager for macOS"
  homepage "https://github.com/kena421/timan"
  url "https://github.com/kena421/timan/archive/refs/tags/v1.0.2.tar.gz"
  version "1.0.2"
  sha256 "1cd76959c2965348bf50f100e904eb19ad760a465e244e3182eee06b4be71ebc"
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
