FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o launchguard-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/launchguard-api .

EXPOSE 8080

CMD ["./launchguard-api"]