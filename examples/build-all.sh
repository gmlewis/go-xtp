#!/bin/bash -e
for i in $(echo */); do
    pushd $i && ./build.sh "$@" && popd
done
