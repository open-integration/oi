
#!/bin/bash

set -e
echo $GOOGLE_SERVICE_ACCOUNT_B64 | base64 -d - > /tmp/service-account.json

echo "Authentication with GCS project $GCS_PROJECT"
cat /tmp/service-account.json
gcloud auth activate-service-account --key-file /tmp/service-account.json
echo "Authenticated!"

FILES=$(ls -ls dist/catalog | awk 'NR>1' | awk '{print $10}')
echo "Files to be uploaded:"
echo "$FILES"
gsutil mv dist/catalog/* $GCS_BUCKET
