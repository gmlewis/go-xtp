#!/bin/bash -e
echo GENERATING user for Go
go run ../../cmd/xtp2code/main.go \
    -lang=go \
    -pkg=user \
    -host=go-host \
    -plugin=go-plugin \
    -types=go-types \
    -yaml=schema.yaml \
    "$@"
echo GENERATING user for MoonBit
go run ../../cmd/xtp2code/main.go \
    -lang=mbt \
    -pkg=user \
    -host=mbt-host \
    -plugin=mbt-plugin \
    -types=mbt-types \
    -yaml=schema.yaml \
    "$@"

for i in $(echo */build.sh); do
    dirname=$(pwd)/${i%"build.sh"}
    echo ENTER ${dirname}
    pushd ${i%"build.sh"} && ( ./build.sh || echo "BUILD FAILED in ${dirname}" ) && popd
    echo LEAVE ${dirname}
done
