package="github.com/aboxofsox/optics"
package_split=(${package//\//})
package_name=${packaeg_split[-1]}

platforms=("windows/amd64", "darwin/amd64", "linux/amd64")

if [! -d "./bin"]; then
    mkdir "bin"
fi

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\//})
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name=$package_name'-'$GOOS'-'$GOARCH

    if [ $GOOS="windows"]; then
        output_name+='.exe'
        env GOOS=$GOOS GOARCH=$GOARCH && go build -tag=windows -o './bin'+$output_name $package
    fi

    env GOOS=$GOOS GOARCH=$GOARCH && go build -tag=linux,darwin -o '/bin/'+$output_name $package
done
