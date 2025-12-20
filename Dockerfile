# ---- Build stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Instalar swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generar docs
RUN swag init

# Compilar
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ---- Runtime stage ----
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
