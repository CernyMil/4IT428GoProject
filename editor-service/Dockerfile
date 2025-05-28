FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o editor-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/editor-service .
RUN cp firebase-admin-sdk.json .

EXPOSE 8081

CMD ["./editor-service"]