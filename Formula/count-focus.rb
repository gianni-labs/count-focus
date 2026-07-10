class CountFocus < Formula
  desc "Focus timer for the terminal"
  homepage "https://github.com/gianni-labs/count-focus"
  url "https://github.com/gianni-labs/count-focus/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "bd7c2de2bee9c3ededed93a6976757aa350eaae67be5c8d49596f9b2ee315dbf"
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
