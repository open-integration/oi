
#!/bin/bash

set -e

branch=$(git symbolic-ref --short HEAD)
commit=$(git rev-parse --short HEAD)
slug="open-integration/core"

echo $branch > ./branch.var
echo $commit > ./commit.var
echo $slug > ./slug.var