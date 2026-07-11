class CountFocus < Formula
  desc "Focus timer for the terminal"
  homepage "https://github.com/gianni-labs/count-focus"
  url "https://github.com/gianni-labs/count-focus/archive/refs/tags/v0.3.1.tar.gz"
  sha256 "e43cd5d50fc49bdc017435412ace08172d646144b86a00cd022e1e774d4b5f1e"
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
