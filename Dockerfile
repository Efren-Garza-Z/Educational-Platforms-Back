# Paso 1: Construir el binario de Go
FROM golang:1.21-alpine AS builder

# Instalar git y certificados (necesarios para llamadas externas como Gemini)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar el proyecto (asegúrate de que el archivo con la func main esté en la raíz o ajusta la ruta)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Paso 2: Crear la imagen ligera de ejecución
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde el constructor
COPY --from=builder /app/main .

# Exponer el puerto que usa Gin (normalmente 8080)
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]