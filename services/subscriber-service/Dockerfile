FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o /app/subscriber-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/subscriber-service/main/api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscriber-api .

EXPOSE 8083

CMD ["./subscriber-service"]
