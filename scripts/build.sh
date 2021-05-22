
#!/bin/bash

set -e

SERVICES="http github slack"
for s in $SERVICES
do  
    V=$(yq e '.version' catalog/services/$s/service.yaml)
    NAME=$(echo $s | tr '[a-z]' '[A-Z]')
    FULL_NAME=$(echo "$NAME""_SERVICE_VERSION")
    export $FULL_NAME=$V
done

goreleaser build --skip-validate --rm-dist --single-target