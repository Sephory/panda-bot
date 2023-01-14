FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . ./

RUN apk update && apk add build-base

RUN go mod download
RUN go build -o panda ./cmd/panda/main.go

FROM alpine:3.17

WORKDIR /app

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build /app/panda ./panda

ENTRYPOINT ["/app/panda", "serve"]
