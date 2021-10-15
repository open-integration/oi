#/bin/bash

# Github
quicktype -l go \
    -s schema \
    --package types \
    --src catalog/services/github/configs/types/types.json \
    --src catalog/services/github/configs/endpoints/getissuecomments/arguments.json \
    --src catalog/services/github/configs/endpoints/getissuecomments/returns.json \
    --src catalog/services/github/configs/endpoints/issuesearch/arguments.json \
    --src catalog/services/github/configs/endpoints/issuesearch/returns.json \
    -o catalog/services/github/types/types.go

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
    -o catalog/services/trello/types/types.go

# Airtable
quicktype -l go \
    -s schema \
    --package types \
    --src catalog/services/airtable/configs/endpoints/addrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/addrecords/returns.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/arguments.json \
    --src catalog/services/airtable/configs/endpoints/getrecords/returns.json \
    -o catalog/services/airtable/types/types.go


# Google Calendar
quicktype -l go \
    -s schema \
    --package types \
    --src catalog/services/google-calendar/configs/endpoints/getEvents/arguments.json \
    --src catalog/services/google-calendar/configs/endpoints/getEvents/returns.json \
    -o catalog/services/google-calendar/types/types.go