# Timezone Rest API

## Description

Small webserver that can deliver timezone information for geo coordinates. For example, you can request this resource:

```
http://localhost:8080/api?lng=52.517932&lat=13.402992
```

And, the service will respond with a json:

```json
{
  "name": "Europe/Berlin",
  "id": "CEST",
  "offset": 7200
}
```

`name` is the common name for the timezone. `id` is the short identifier. `offset` is the difference in seconds to UTC.

## How to use

This project based on the new go modules. So you only need `go >= 1.11`. Change in project folder and call `go build` to build the binaries or `go run main.go` to start the application directly.

### Docker

```
docker run -p8080:8080 racemap/tz-service
```

If you want to build the container yourself, change to the project folder and run `docker build -t tz-service .`. To start the container after build run `docker run -p 8080:8080 tz-service`.

## Thanks to

This project use the work of two awesome projects. [evanoberholster/timezoneLookup](https://github.com/evanoberholster/timezoneLookup) is go library that find the timezone for a coordinate. This project use the shapefile of this project [evansiroky/timezone-boundary-builder](https://github.com/evansiroky/timezone-boundary-builder).

## TODO

- [x] Better documentation
- [ ] Tests
