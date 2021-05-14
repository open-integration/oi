
#!/bin/bash

set -e

go mod tidy
if [[ -z $(git status -s) ]]
then
    echo "Directory is clean"
else
    echo "Working directory is now clean"
    git status -s
    exit 1
fi