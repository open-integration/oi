# OPENC - a pipeline execution engine

## stop using YAMLs

## This is not production ready stuff

## Concepts
* A compiled, binray pipeline
* State - Openc holds the state of all the tasks
* Service - a standalone binary exposing API over http2 (gRPC) that Openc can trigger, called endpoint.
* Task - a request from a service to run some logic.

* Endpoint of a service defined by 2 files of JSON schema, `arguments.json` and `returns.json`, Openc will enforce the arguments given by a task and the output created to match the schema.

## Example
* An example of pipeline can be found [here](https://github.com/olegsu/trello-sync). The pipeline scans my Trello board and update a rows in my Google Spreadsheet, once a card been moved to "Done" list, it will be archived from the board.
