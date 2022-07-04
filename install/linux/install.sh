# simple install script. requires sudo
env GOOS='linux' GOARCH='amd64' && go build -tag=linux -o './optics'
mv './optics' /usr/bin
