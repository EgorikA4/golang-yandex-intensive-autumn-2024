FROM golang:1.24-alpine3.21 AS builder

WORKDIR /code

COPY . .

RUN go build -o orchestrator ./cmd/orchestrator
RUN go build -o agent ./cmd/agent

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /code/orchestrator .
COPY --from=builder /code/agent .
COPY --from=builder /code/.env .
