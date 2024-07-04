#!/bin/bash -e
echo PUSH and BIND user for Go
pushd go-plugin
xtp plugin push || echo "xtp push FAILED"
xtp plugin bind || echo "xtp bind FAILED"
popd

echo PUSH and BIND user for MoonBit
pushd mbt-plugin
xtp plugin push || echo "xtp push FAILED"
xtp plugin bind || echo "xtp bind FAILED"
popd
