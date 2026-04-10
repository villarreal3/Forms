FROM docker.io/library/golang:1.21-alpine AS builder

WORKDIR /app

# Instala dependencias del sistema necesarias para compilar
RUN apk add --no-cache git build-base

# Copia los archivos de módulos primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código (incl. docs/ para go:embed de OpenAPI + Swagger UI)
COPY . .

# Compila el binario (spec embebido; no hace falta copiar docs/ en la imagen final)
RUN go build -o app .

# Imagen final ligera
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache mysql-client ca-certificates tzdata && \
    update-ca-certificates

# Copia el binario compilado desde la imagen de build
COPY --from=builder /app/app /app/app

# Variables de entorno típicas para conexión a MySQL
ENV DB_HOST=mysql \
    DB_PORT=3306 \
    DB_USER=root \
    DB_PASSWORD=secret \
    DB_NAME=app_db \
    SERVER_PORT=8080 \
    EMAILS_DIR=src/emails \
    IMAGES_DIR=src/img

# Expón el puerto en el que corre tu API (ajusta si es distinto)
EXPOSE 8080

# Comando de arranque
CMD ["/app/app"]

