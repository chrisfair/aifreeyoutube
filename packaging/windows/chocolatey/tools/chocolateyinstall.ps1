# ==============================================================================
# Chocolatey Installation Script for AI-Free YouTube
# ==============================================================================
$ErrorActionPreference = 'Stop'

$packageName = 'aifreeyoutube'
$toolsDir     = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$installDir   = Join-Path $([System.Environment]::GetFolderPath('LocalApplicationData')) 'Programs\aifreeyoutube'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $installDir
  url64         = "https://github.com/christopherpatrickfair/aifreeyoutube/releases/download/v1.0.0/aifreeyoutube-v1.0.0-windows-amd64.zip"
  checksum64    = "REPLACE_WITH_SHA256"
  checksumType64= 'sha256'
}

# Download and unpack zip release bundle
Install-ChocolateyZipPackage @packageArgs

# Ensure destination directory exists
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

$exePath = Join-Path $installDir 'aifreeyoutube.exe'

# Create Desktop and Start Menu Shortcuts
$desktopPath   = [System.Environment]::GetFolderPath('Desktop')
$startMenuPath = [System.Environment]::GetFolderPath('Programs')

Install-ChocolateyShortcut -ShortcutFilePath (Join-Path $desktopPath 'AI-Free YouTube.lnk') -TargetPath $exePath -IconLocation $exePath
Install-ChocolateyShortcut -ShortcutFilePath (Join-Path $startMenuPath 'AI-Free YouTube.lnk') -TargetPath $exePath -IconLocation $exePath

# Add installation directory to User PATH
Install-ChocolateyPath $installDir 'User'
Write-Host "AI-Free YouTube successfully installed to $installDir"
