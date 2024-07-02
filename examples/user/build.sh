#!/bin/bash -e
go run ../../cmd/xtp2code/main.go \
    -lang=go \
    -pkg=user \
    -host=go-host \
    -plugin=go-plugin \
    -types=go-types \
    -yaml=schema.yaml
go run ../../cmd/xtp2code/main.go \
    -lang=mbt \
    -pkg=user \
    -host=mbt-host \
    -plugin=mbt-plugin \
    -types=mbt-types \
    -yaml=schema.yaml

for i in $(echo */build.sh); do
    pushd ${i%"build.sh"} && ./build.sh && popd
done
