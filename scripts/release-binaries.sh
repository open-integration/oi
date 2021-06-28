
#!/bin/bash

set -e

SERVICES="http github slack"
for s in $SERVICES
do  
    V=$(cat catalog/services/$s/VERSION)
    NAME=$(echo $s | tr '[a-z]' '[A-Z]')
    FULL_NAME=$(echo "$NAME""_SERVICE_VERSION")
    export $FULL_NAME=$V
done

goreleaser release --rm-dist