
#!/bin/bash

set -e

SERVICES="http github slack trello airtable google-calendar"
for s in $SERVICES
do  
    V=$(cat catalog/services/$s/VERSION)
    NAME=$(echo $s | tr '[a-z]' '[A-Z]' | sed s/-/_/g)
    FULL_NAME=$(echo "$NAME""_SERVICE_VERSION")
    export $FULL_NAME=$V
done

goreleaser build --skip-validate --rm-dist --single-target --snapshot