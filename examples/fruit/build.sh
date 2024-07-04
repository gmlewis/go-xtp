#!/bin/bash -e
echo GENERATING fruit for Go
go run ../../cmd/xtp2code/main.go \
    -lang=go \
    -pkg=fruit \
    -host=go-host \
    -plugin=go-plugin \
    -types=go-types \
    -yaml=schema.yaml \
    "$@"
echo GENERATING fruit for MoonBit
go run ../../cmd/xtp2code/main.go \
    -lang=mbt \
    -pkg=fruit \
    -host=mbt-host \
    -plugin=mbt-plugin \
    -types=mbt-types \
    -yaml=schema.yaml \
    "$@"

for i in $(echo */build.sh); do
    echo ENTER ${i%"build.sh"}
    pushd ${i%"build.sh"} && ./build.sh && popd
    echo LEAVE ${i%"build.sh"}
done
