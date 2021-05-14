
#!/bin/bash

set -e

version=$(cat VERSION)
echo "Releasing version $version"
git tag v${version}
git push --tags 