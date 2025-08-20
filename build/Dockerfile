FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
COPY ../cmd ./cmd
COPY ../internal ./internal

RUN go mod download

RUN GOOS=linux go build -o wb-lo-task ./cmd

FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /app/wb-lo-task .
COPY ../config.yaml ./config.yaml
COPY ../web ./web

EXPOSE 8080

CMD ["./wb-lo-task"]