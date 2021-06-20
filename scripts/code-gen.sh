#/bin/bash

# TODO: get all the data from service.yaml files
mkdir -p catalog/services/github/endpoints/issuesearch
quicktype -o catalog/services/github/endpoints/issuesearch/arguments.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/arguments.json --package issuesearch -t IssueSearchArguments
quicktype -o catalog/services/github/endpoints/issuesearch/returns.go -l go -s schema --src catalog/services/github/configs/endpoints/issuesearch/returns.json --package issuesearch -t IssueSearchReturns

mkdir -p catalog/services/github/endpoints/getissuecomments
quicktype -o catalog/services/github/endpoints/getissuecomments/arguments.go -l go -s schema --src catalog/services/github/configs/endpoints/getissuecomments/arguments.json --package getissuecomments -t GeIssueCommentArguments
quicktype -o catalog/services/github/endpoints/getissuecomments/returns.go -l go -s schema --src catalog/services/github/configs/endpoints/getissuecomments/returns.json --package getissuecomments -t GetIssueCommentsReturns

mkdir -p catalog/services/trello/endpoints/getcards
quicktype -o catalog/services/trello/endpoints/getcards/returns.go -l go -s schema --src catalog/services/trello/configs/endpoints/getcards/returns.json --package getcards -t GetcardsReturns

# Airtable
quicktype -l go \
    -s schema \
    --src catalog/services/airtable/configs/endpoints/addrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/addrecords/returns.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/returns.json \
    --package types > catalog/services/airtable/types/types.go
