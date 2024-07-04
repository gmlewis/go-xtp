#!/bin/bash -e
echo PUSH and BIND fruit for Go
pushd go-plugin
xtp plugin push || echo "xtp push FAILED"
xtp plugin bind || echo "xtp bind FAILED"
popd

echo PUSH and BIND fruit for MoonBit
pushd mbt-plugin
xtp plugin push || echo "xtp push FAILED"
xtp plugin bind || echo "xtp bind FAILED"
popd
