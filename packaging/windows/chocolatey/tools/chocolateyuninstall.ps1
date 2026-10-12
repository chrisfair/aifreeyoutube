# ==============================================================================
# Chocolatey Uninstallation Script for AI-Free YouTube
# ==============================================================================
$ErrorActionPreference = 'Continue'

$packageName = 'aifreeyoutube'
$installDir  = Join-Path $([System.Environment]::GetFolderPath('LocalApplicationData')) 'Programs\aifreeyoutube'

# Remove shortcuts
$desktopPath   = [System.Environment]::GetFolderPath('Desktop')
$startMenuPath = [System.Environment]::GetFolderPath('Programs')

Remove-Item (Join-Path $desktopPath 'AI-Free YouTube.lnk') -ErrorAction SilentlyContinue
Remove-Item (Join-Path $startMenuPath 'AI-Free YouTube.lnk') -ErrorAction SilentlyContinue

# Remove PATH entry
Uninstall-ChocolateyPath $installDir 'User'

# Clean up binaries
if (Test-Path $installDir) {
    Remove-Item $installDir -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host "AI-Free YouTube uninstalled successfully."
