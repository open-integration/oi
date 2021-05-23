# Open Integration

Open-Integration is a pipeline framework that allows to write (golang) powerful pipelines, then the pipeline is compiled into single binary to be executed elsewhere.

## Why?
A pipeline, wether its CI, CD, ETL and any other binding between multiple API's is most of the time 2 parts entity:
1. The pipeline engine
2. The pipeline spec

Its hard to test
Its hard to run locally
Required knowledge of providers pipeline spec

Well besicaly its once again the big debate code vs config
