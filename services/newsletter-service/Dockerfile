FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o newsletter-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/newsletter-service .
#COPY firebase-admin-sdk.json . 

EXPOSE 8082

CMD ["./newsletter-service"]