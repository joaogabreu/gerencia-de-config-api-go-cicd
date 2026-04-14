# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-go ./main.go

FROM alpine:3.21
WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/api-go /app/api-go

EXPOSE 3000
USER appuser

CMD ["/app/api-go"]
