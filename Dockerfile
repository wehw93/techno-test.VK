# Этап сборки
FROM golang:1.22.2 AS builder

WORKDIR /app


RUN apt-get update && \
    apt-get install -y gcc libssl-dev pkg-config && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download


COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY tarantool/ ./tarantool/


RUN CGO_ENABLED=1 GOOS=linux go build -o voting-bot ./cmd/voting-bot


FROM debian:stable-slim


RUN apt-get update && \
    apt-get install -y libssl-dev ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY --from=builder /app/voting-bot .

COPY --from=builder /app/tarantool/init.lua /opt/tarantool/


COPY local.env .

CMD ["./voting-bot"]