
#!/bin/bash

set -e

fmtcmd="gofmt -l . | wc -l"
if [ $(eval $fmtcmd) -gt 0 ]
then
    echo "cmd: \"$fmtcmd\" failed"
    exit 1
fi