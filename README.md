# Timezone Rest Api

## Description:

Small webserver, that can deliver Timezone Informations for coordinates. For example you can request this resource: 

```
http://localhost:8080/api?lng=52.517932&lat=13.402992"
```

The service response with a json:

```json
{
    "name": "Europe/Berlin",
    "id": "CEST",
    "offset": 7200
}
```

## Thanks to:

This project use the work of two awesome projects. [evanoberholster/timezoneLookup](https://github.com/evanoberholster/timezoneLookup) is go library that find the timezone for a coordinate. This project use the shapefile of this project [evansiroky/timezone-boundary-builder](https://github.com/evansiroky/timezone-boundary-builder). 

## How to use:

This project based on the new go modules. So you only need `go >= 1.11`. Change in project folder and call `go build` to build the binaries or `go run main.go` to start the application directly.

## Docker:

Or you use docker. Then you can change in the project folder and call `docker build -t tz-service .` and after the build `docker run -p 8080:8080 tz-service`.

## TODO:

[ ] Better documentation
[ ] Tests

