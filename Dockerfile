# Cambiamos 1.21-alpine por alpine (que es la última versión estable)
FROM golang:alpine AS builder

# Instalamos certificados y herramientas necesarias
RUN apk update && apk add --no-cache ca-certificates git

WORKDIR /app

# Copiamos dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos
# Si tu main.go está en la raíz, deja el punto. 
# Si está en /cmd, usa ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- Etapa de ejecución ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
