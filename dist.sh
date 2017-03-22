#!/bin/bash

# build binary distributions for linux/amd64 and darwin/amd64
set -e 

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"

echo "... running tests"
go test|| exit 1

arch=$(go env GOARCH)
version=$(cat $DIR/version.go | grep "const VERSION" | awk '{print $NF}' | sed 's/"//g')
goversion=$(go version | awk '{print $3}')

for os in linux darwin; do
    echo "... building v$version for $os/$arch"
    BUILD=$(mktemp -d -t json2csv)
    TARGET="json2csv-$version.$os-$arch.$goversion"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build
    mkdir -p $BUILD/$TARGET
    cp json2csv $BUILD/$TARGET/json2csv
    pushd $BUILD >/dev/null
    tar czvf $TARGET.tar.gz $TARGET
    if [ -e $DIR/dist/$TARGET.tar.gz ]; then
        echo "... WARNING overwriting dist/$TARGET.tar.gz"
    fi
    mv $TARGET.tar.gz $DIR/dist
    echo "... built dist/$TARGET.tar.gz"
    popd >/dev/null
done