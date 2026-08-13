FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /weather-api ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /weather-api /weather-api

EXPOSE 8080

ENTRYPOINT ["/weather-api"]
