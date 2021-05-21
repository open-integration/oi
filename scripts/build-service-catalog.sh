
#!/bin/bash

set -e

build () {
    NAME=$1
    VERSION=$2
    dest="./dist/catalog/$NAME-$VERSION"
    pkg="cmd/catalog/$NAME/main.go"
    
    echo "Building $NAME-$VERSION service for darwin and amd64"
    GOOS=darwin GOARCH=amd64 go build -o $dest-darwin-amd64 $pkg 

    echo "Building $NAME-$VERSION service for linux and amd64"
    GOOS=linux GOARCH=amd64 go build -o $dest-linux-amd64 $pkg

    echo "Building $NAME-$VERSION service for linux and 386"
    GOOS=linux GOARCH=386 go build -o $dest-linux-386 $pkg
}

SERVICES="http github"

echo ""

for s in $SERVICES
do  
    V=$(yq e '.version' catalog/services/$s/service.yaml)
    echo "Building service $s ($V)"
    build $s $V
done