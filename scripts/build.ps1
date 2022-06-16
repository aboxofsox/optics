$Package = "github.com/aboxofsox/optics"
$PackageSplit = $Package.Split("/")
$PackageName = $PackageSplit[-1]

$Platforms = "windows/amd64", "darwin/amd64", "linux/amd64"

ForEach ($Platform in $Platforms) {
    $Emoji = "🖥️"
    $PlatformSplit = $Platform.Split("/")
    If ($PlatformSplit[0] -eq "windows") {
        $Emoji = "🪟"
    }
    If ($PlatformSplit[0] -eq "darwin") {
        $Emoji = "🍎"
    }
    If ($PlatformSplit[0] -eq "linux") {
        $Emoji = "🐧"
    }

    Write-Host "Building for $Platform $Emoji" -ForegroundColor Cyan
    $GOOS = $PlatformSplit[0]
    $GOARCH = $PlatformSplit[1]
    $Env:GOOS = $GOOS
    $Env:GOARCh = $GOARCH
    
    $OutName = "./bin/$PackageName-$GOOS-$GOARCH"

    If ($GOOS -eq "windows") {
        $OutName += ".exe"
    }

    cmd.exe /c "go build -o $OutName $Package"
}

If ((Test-Path "./bin/optics.sum") -ne $True) {
    New-Item -Path "./bin/optics.sum" -ItemType "file"
}
$BinContent = Get-ChildItem -Path "./bin"
For ($i = 0; $i -lt $BinContent.Length; $i++) {
    $Sum = Get-FileHash $BinContent[$i] | Select-Object -Property Hash -ExpandProperty Hash
    $Name = $BinContent[$i].Name
    Add-Content -Path "./bin/optics.sum" -Value "$Name : $Sum"
}

Write-Host "Build complete ✅" -ForegroundColor Green

