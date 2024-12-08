FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -ldflags="-s -w" -o /app/service ./cmd/app/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/service /app/service

EXPOSE 8080

CMD ["./service"]