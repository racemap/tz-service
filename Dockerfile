# build stage
FROM golang:alpine3.13 AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

RUN /app/tz-service init

# final stage
FROM alpine:3.13 AS app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/tz-service /app/
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/

RUN apk add --update --no-cache curl

HEALTHCHECK CMD curl --fail "http://localhost:8080/api?lng=52.517932&lat=13.402992" || exit 1

EXPOSE 8080

WORKDIR /app

ENTRYPOINT ["/app/tz-service"]

FROM app AS app_with_assets
COPY --from=builder /app/assets/timezone.snap.json /app/assets/
