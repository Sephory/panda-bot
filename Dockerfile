FROM golang:1.19-alpine

WORKDIR /app

COPY . ./

RUN apk update && apk add build-base ca-certificates && rm -rf /var/cache/apk/*

RUN go mod download

ENTRYPOINT ["go", "run", "./cmd/pandabot/main.go", "serve", "--http=0.0.0.0:9000"]
