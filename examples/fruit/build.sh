#!/bin/bash -e
# go run ../../cmd/xtp2code/main.go \
#     -lang=go \
#     -host=go-host \
#     -plugin=go-plugin \
#     -types=go-types/fruit.go \
#     -yaml=schema.yaml
# go run ../../cmd/xtp2code/main.go \
#     -lang=mbt \
#     -host=mbt-host \
#     -plugin=mbt-plugin \
#     -types=mbt-types/fruit.mbt \
#     -yaml=schema.yaml

pushd go-plugin && ./build.sh && popd
