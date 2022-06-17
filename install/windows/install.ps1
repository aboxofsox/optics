# Simple Windows install script

$InstallPath = "C:\\Users\\$($Env:USERNAME)\\optics"

If ((Test-Path -Path $InstallPath) -ne $True) {
    New-Item -Path $InstallPath -Name 'optics' -ItemType 'directory' | Out-Null
}

# Create a temporary file to hold the download.
$Temp = New-TemporaryFile

# Download the latest build of Optics for Windows
Invoke-WebRequest -Uri 'https://github.com/aboxofsox/optics/releases/download/1.1.0/optics-windows-amd64.exe' -Outfile $Temp

# Validate download.
# If validation fails, exit the script.
If ((Test-Path -Path $Temp) -ne $True) {
    Write-Host 'Unable to verify download.' -ForegroundColor Red
    Exit
}

# Copy the temporary file to C:\ProgramFiles\optics
Move-Item -Path $Temp -Destination $InstallPath
Rename-Item -Path "$InstallPath\\$($Temp.Name)" -NewName "$InstallPath\\optics.exe"

# Add the directory to environment path
$Env:PATH += ";$InstallPath"

Write-Host 'Installation complete ✅' -ForegroundColor Green
Write-Host "Type 'optics init' to get started." -Foreground Gray