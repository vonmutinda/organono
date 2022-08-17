# build
FROM docker.io/golang:1.18-alpine as builder

WORKDIR /app

COPY . /app/

RUN \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o organono-api cmd/main.go

# run
FROM docker.io/alpine:3.14

WORKDIR /app

RUN \
    apk update

COPY \
    --from=builder /app/organono-api /app/.env /app/

EXPOSE 80

ENTRYPOINT ["./organono-api"]
