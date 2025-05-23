FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o subscriber-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscriber-service .
COPY goproject-4c949-firebase-adminsdk-fbsvc-564d6d7e63.json . 

EXPOSE 8083

CMD ["./subscriber-service"]
