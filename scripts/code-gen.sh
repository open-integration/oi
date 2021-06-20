#/bin/bash

# TODO: get all the data from service.yaml files
mkdir -p catalog/services/github/endpoints/issuesearch
quicktype -o catalog/services/github/endpoints/issuesearch/arguments.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/arguments.json --package issuesearch -t IssueSearchArguments
quicktype -o catalog/services/github/endpoints/issuesearch/returns.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/returns.json --package issuesearch -t IssueSearchReturns

mkdir -p catalog/services/github/endpoints/getissuecomments
quicktype -o catalog/services/github/endpoints/getissuecomments/arguments.go -l go -s schema --src catalog/services/github/configs/endpoints/getissuecomments/arguments.json --package getissuecomments -t GeIssueCommentArguments
quicktype -o catalog/services/github/endpoints/getissuecomments/returns.go -l go -s schema --src catalog/services/github/configs/endpoints/getissuecomments/returns.json --package getissuecomments -t GetIssueCommentsReturns

# Trello
quicktype -l go \
    -s schema \
    --package types \
    --src catalog/services/trello/configs/endpoints/getcards/arguments.json \
    --src catalog/services/trello/configs/endpoints/getcards/returns.json \
    --src catalog/services/trello/configs/endpoints/archivecard/arguments.json \
    --src catalog/services/trello/configs/endpoints/archivecard/returns.json \
    --src catalog/services/trello/configs/endpoints/addcard/arguments.json \
    --src catalog/services/trello/configs/endpoints/addcard/returns.json \
    > catalog/services/trello/types/types.go

# Airtable
quicktype -l go \
    -s schema \
    --package types \
    --src catalog/services/airtable/configs/endpoints/addrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/addrecords/returns.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/returns.json \
    > catalog/services/airtable/types/types.go
