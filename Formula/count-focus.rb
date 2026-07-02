class CountFocus < Formula
  desc "Focus timer for the terminal"
  homepage "https://github.com/gianni-labs/count-focus"
  url "https://codeload.github.com/gianni-labs/count-focus/tar.gz/refs/tags/v0.1.0"
  sha256 "ea88f7ff34cf690ccab6695c724c7f27a5c4b4869a6d21d0cbed8c2a331c23ee"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-trimpath", "-ldflags", "-s -w", "-o", bin/"count-focus", "."
  end

  test do
    assert_match "Usage:", shell_output("#{bin}/count-focus --help")
  end
end
