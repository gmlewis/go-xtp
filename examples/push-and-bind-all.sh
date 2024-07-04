#!/bin/bash -e
for i in $(echo */); do
    echo ENTER $i
    pushd $i && ./push-and-bind.sh "$@" && popd
    echo LEAVE $i
done
