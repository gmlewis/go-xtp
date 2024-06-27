#!/bin/bash -e
go run ../../cmd/xtp2code/main.go \
    -lang=go \
    -host=go-host/main.go \
    -plugin=go-plugin/main.go \
    -types=go-types/fruit.go \
    -yaml=schema.yaml
go run ../../cmd/xtp2code/main.go \
    -lang=mbt \
    -host=mbt-host/main.mbt \
    -plugin=mbt-plugin/main.mbt \
    -types=mbt-types/fruit.mbt \
    -yaml=schema.yaml
