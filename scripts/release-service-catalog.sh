
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

SERVICES="http"

echo ""

for s in $SERVICES
do  
    V=$(yq e '.version' catalog/services/$s/service.yaml)
    echo "Building service $s ($V)"
    build $s $V
done


echo $GOOGLE_SERVICE_ACCOUNT_B64 | base64 -d - > /tmp/service-account.json

echo "Authentication with GCS project $GCS_PROJECT"
cat /tmp/service-account.json
gcloud auth activate-service-account --key-file /tmp/service-account.json
echo "Authenticated!"

FILES=$(ls -ls dist/catalog | awk 'NR>1' | awk '{print $10}')
echo "Files to be uploaded:"
echo "$FILES"
gsutil mv dist/catalog/* $GCS_BUCKET
