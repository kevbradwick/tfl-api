TfL API
=======

A HTTP API hooked up to Station data prided by [Transport for London](https://api-portal.tfl.gov.uk).

This project provides a CLI ingest tool to consume the station facilities xml feed to
a MongoDB instance. It also builds a HTTP API (`api.go`) to query the data.

### Installation

Load the data into Mongo using the ingest tool

    go run ingest/ingest.go --input=data/stationsfacilities.xml
    
It will connect to the instance running on `localhost:27071` by default. The
options `--dbhost` and `--dbport` can be used to override these values.
        
Then run the API

    go run api/api.go
    
### Handlers

#### GET /stations/search

Search all the stations. The following query parameters can be used to filter
the results;

* lines
* zones
* name

You can also sort by distance (nearest first) from a point using the parameter
`near` that accepts a *longitude*, *latitude*, for example

    /stations/search?near=-0.2152,52.123
    
#### GET /stations/station/:id

Get a single station e.g.

    /stations/station/1000002
