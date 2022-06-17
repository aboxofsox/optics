# Script to install Optics
# Because it makes a permanent env path change, it requies elevation.

# Set the version
$Version = '1.1.0'

# Check if user is admin.
# This is required for the permanent path change.
$IsAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
If ($IsAdmin -ne $True) {
    Write-Host 'Elevation required for permanent path change.' -ForegroundColor Red
    Exit
}

# Set the installation directory.
# Optics will be installed within the user scope.
$InstallPath = "C:\Users\$($Env:USERNAME)\Optics"

If ((Test-Path -Path $InstallPath) -ne $True) {
    New-Item -Path $InstallPath -Name 'optics' -ItemType 'directory' | Out-Null
}

# Create a temporary file to hold the download
$Temp = New-TemporaryFile

# Download the latest build of Optics for Windows
Invoke-WebRequest -Uri "https://github.com/aboxofsox/optics/releases/download/$Version/optics-windows-amd64.exe" -OutFile $Temp

# Validate that the file exists.
If ($Temp.Exists -ne $True) {
    Write-Host 'Unable to validate download.' -ForegroundColor Red
    Exit
}

# Copy the temporary file to C:\Users\{username}\optics
Move-Item -Path $Temp -Destination $InstallPath
Rename-Item -Path "$InstallPath\$($Temp.Name)" -NewName "$InstallPath\optics.exe"

# Add the path to env paths
# This requires elevation since it's a registry key change
$RegPath = 'Registry::HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment'

If ((Test-Path -Path $RegPath) -ne $True) {
    Write-Host 'Registry key does not exist.' -ForegroundColor Red
    Exit
}

$OriginalPath = (Get-ItemProperty -Path $RegPath -Name PATH).Path
$AddPath = "$OriginalPath;$InstallPath"
Set-ItemProperty -Path $RegPath -Name PATH -Value $AddPath

# Validate the environment path
$TestEnvPath = (Get-ItemProperty -Path $RegPath -Name Path).Path
$TestEnvPathSplit = $TestEnvPath.Split(';')
If ($TestEnvPathSplit.Contains($InstallPath) -ne $True) {
    Write-Host 'Path change was not successful.' -ForegroundColor Yellow
}

Write-Host 'Installation complete.' -ForegroundColor Green
Write-Host "Type 'optics init' to get started." -ForegroundColor Gray