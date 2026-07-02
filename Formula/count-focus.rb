class CountFocus < Formula
  desc "Focus timer for the terminal"
  homepage "https://github.com/gianni-labs/count-focus"
  url "https://github.com/gianni-labs/count-focus.git", branch: "main"
  version "0.1.0-dev"

  depends_on "go" => :build

  def install
    system "go", "build", "-trimpath", "-ldflags", "-s -w", "-o", bin/"count-focus", "."
  end

  test do
    assert_match "Usage:", shell_output("#{bin}/count-focus --help")
  end
end
