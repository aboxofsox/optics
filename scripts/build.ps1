$Package = "github.com/aboxofsox/optics"
$PackageSplit = $Package.Split("/")
$PackageName = $PackageSplit[-1]

$Platforms = "windows/amd64", "darwin/amd64", "linux/amd64"

$CurrentItem = 0
$PercentComplete = 0

For ($i = 0; $i -lt $Platforms.Length; $i++) {
    Write-Progress -Activity 'Building Optics' -Status "$PercentComplete% Complete:" -PercentComplete $PercentComplete 
    $Emoji = "🖥️"
    $PlatformSplit = $Platforms[$i].Split("/")
    If ($PlatformSplit[0] -eq "windows") {
        $Emoji = "🪟"
    }
    If ($PlatformSplit[0] -eq "darwin") {
        $Emoji = "🍎"
    }
    If ($PlatformSplit[0] -eq "linux") {
        $Emoji = "🐧"
    }

    Write-Host "Building for $($Platforms[$i]) $Emoji" -ForegroundColor Cyan
    $GOOS = $PlatformSplit[0]
    $GOARCH = $PlatformSplit[1]
    $Env:GOOS = $GOOS
    $Env:GOARCh = $GOARCH
    
    $OutName = "./bin/$PackageName-$GOOS-$GOARCH"

    If ($GOOS -eq "windows") {
        $OutName += ".exe"
        cmd.exe /c "go build -tags=windows -o $OutName $Package"
    } Else {
        cmd.exe /c "go build -tags=lnx,dar -o $OutName $Package"
    }
    $CurrentItem++
    $PercentComplete = [int](($CurrentItem / $Platforms.Length) * 100)
} 

Write-Host "Build complete ✅" -ForegroundColor Green

Write-Host ""
$BinContent = Get-ChildItem -Path "./bin"
For ($i = 0; $i -lt $BinContent.Length; $i++) {
    $Sum = Get-FileHash $BinContent[$i] | Select-Object -Property Hash -ExpandProperty Hash
    $Obj = @{}
    $Obj[$BinContent[$i].Name] = $Sum
    $Obj
    
}
Write-Host ""
