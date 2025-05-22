# Použij oficiální Go image jako builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Zkopíruj go mod a sum soubory
COPY go.mod go.sum ./
RUN go mod download

# Zkopíruj zbytek zdrojového kódu
COPY . .

# Sestav binárku
RUN go build -o editor-service ./cmd/api

# Použij malý image pro běh aplikace
FROM alpine:latest

WORKDIR /app

# Zkopíruj binárku z builderu
COPY --from=builder /app/editor-service .

# Nastav port (případně uprav dle potřeby)
EXPOSE 8081

# Spusť službu
CMD ["./editor-service"]