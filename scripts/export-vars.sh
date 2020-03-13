
#!/bin/bash

set -e

branch=$(git symbolic-ref --short HEAD)
commit=$(git rev-parse HEAD)
slug="open-integration/core"
version=$(cat VERSION)

echo $branch > ./branch.var
echo $commit > ./commit.var
echo $slug > ./slug.var
echo $version > ./version.var
