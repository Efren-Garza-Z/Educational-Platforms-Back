# --- PASO 1: Compilación ---
FROM golang:1.21-alpine AS builder

# apk update ayuda a refrescar los repositorios antes de instalar
RUN apk update && apk add --no-cache ca-certificates

WORKDIR /app

# Copiamos solo los archivos de dependencias primero (esto optimiza el cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# IMPORTANTE: Si tu archivo main.go está dentro de una carpeta (ej: cmd/main.go)
# cambia el punto final por la ruta: RUN go build -o main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- PASO 2: Imagen de ejecución ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Traemos el ejecutable desde el builder
COPY --from=builder /app/main .

# Exponemos el puerto estándar
EXPOSE 8080

# Comando para arrancar
CMD ["./main"]
