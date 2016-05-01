TfL API
=======

A HTTP API hooked up to Station data prided by [Transport for London](https://api-portal.tfl.gov.uk).

This project provides a CLI ingest tool to consume the station facilities xml feed to
a MongoDB instance. It also builds a HTTP API (`api.go`) to query the data.
