# ==============================================================================
# Homebrew Formula for AI-Free YouTube
# ==============================================================================
# Allows macOS users to install and manage AI-Free YouTube via Homebrew:
#   brew tap christopherpatrickfair/aifreeyoutube
#   brew install aifreeyoutube
#
# Background service management:
#   brew services start aifreeyoutube
# ==============================================================================

class Aifreeyoutube < Formula
  desc "Authentic human video portal filtering out AI slop on YouTube"
  homepage "https://github.com/christopherpatrickfair/aifreeyoutube"
  version "1.0.0"
  license "GPL-3.0-or-later"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/christopherpatrickfair/aifreeyoutube/releases/download/v1.0.0/aifreeyoutube-v1.0.0-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    else
      url "https://github.com/christopherpatrickfair/aifreeyoutube/releases/download/v1.0.0/aifreeyoutube-v1.0.0-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_SHA256"
    end
  end

  on_linux do
    url "https://github.com/christopherpatrickfair/aifreeyoutube/releases/download/v1.0.0/aifreeyoutube_1.0.0_amd64.deb"
    sha256 "PLACEHOLDER_LINUX_SHA256"
  end

  def install
    bin.install "aifreeyoutube-server"
    
    # If desktop app bundle is present in the archive
    if File.directory?("AI-Free YouTube.app")
      prefix.install "AI-Free YouTube.app"
    end
  end

  def caveats
    <<~EOS
      AI-Free YouTube is installed!
      
      To run the web app in your browser:
        aifreeyoutube-server
      Then open: http://localhost:8484
      
      To start it as a continuous background service:
        brew services start aifreeyoutube
    EOS
  end

  service do
    run [opt_bin/"aifreeyoutube-server"]
    keep_alive true
    log_path var/"log/aifreeyoutube.log"
    error_log_path var/"log/aifreeyoutube.err.log"
    environment_variables PORT: "8484"
  end

  test do
    system "#{bin}/aifreeyoutube-server", "--help" rescue true
  end
end
