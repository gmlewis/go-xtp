#!/bin/bash -e
for i in $(echo */); do
    echo ENTER $i
    pushd $i && ./build.sh "$@" && popd
    echo LEAVE $i
done
