class CountFocus < Formula
  desc "Focus timer for the terminal"
  homepage "https://github.com/gianni-labs/count-focus"
  url "https://github.com/gianni-labs/count-focus/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "1f30e0a6cbb7b7c1b5c8ea382c9a0265f3b33c1fe2b5fe36c2b6dbdd2b6f833f"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=v#{version}")
  end

  test do
    assert_match "count-focus <duration>", shell_output("#{bin}/count-focus --help")
    assert_match "count-focus v#{version}", shell_output("#{bin}/count-focus --version")
  end
end
