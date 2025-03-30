
FROM golang:1.22-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


RUN apk add --no-cache openssl-dev gcc musl-dev


COPY . .


RUN CGO_ENABLED=1 GOOS=linux go build -o voting-server ./cmd/server
