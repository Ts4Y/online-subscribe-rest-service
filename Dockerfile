FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./ ./
RUN go build -o ./subscribes ./cmd/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/subscribes .

EXPOSE 8080

ENTRYPOINT [ "/app/subscribes" ]
