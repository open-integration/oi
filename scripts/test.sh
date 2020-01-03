
#!/bin/bash

set -e
set -o pipefail

rm -rf .cover/ .test/
mkdir .cover/ .test/
trap "rm -rf .test/" EXIT

for pkg in `go list ./... | grep -v /vendor/`; do
    go test -v -covermode=atomic \
        -coverprofile=".cover/$(echo $pkg | sed 's/\//_/g').cover.out" $pkg
done
code=$?
go tool cover -html=cp.out -o coverage.html
echo "go test cmd exited with code $code"
exit $code