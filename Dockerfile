FROM golang:alpine AS builder

WORKDIR /code

ADD go.mod .

COPY . .

RUN go build -o orchestrator ./cmd/orchestrator
RUN go build -o agent ./cmd/agent

FROM alpine

WORKDIR /app

COPY --from=builder /code/orchestrator .
COPY --from=builder /code/agent .
COPY --from=builder /code/.env .

EXPOSE 8000

CMD ["./orchestrator", "&", "./agent"]
