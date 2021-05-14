
#!/bin/bash

set -e

version=$(cat VERSION)
echo "Releasing version $version"
fqrn="v$version"
git checkout -b release-$fqrn
git tag $fqrn
git push --tags 