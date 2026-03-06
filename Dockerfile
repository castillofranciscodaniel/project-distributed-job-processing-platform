# ETAPA 1: Compilación
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiamos los archivos de módulos primero para aprovechar cache
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# ETAPA 2: Ejecución (Imagen final ligera)
FROM alpine:latest

# Instalamos certificados (CRUCIAL para hablar con S3) y tzdata
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Copiamos solo el binario desde el builder
COPY --from=builder /app/main .

# Por seguridad en ECS: Creamos un usuario no-root
RUN adduser -D appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["./main"]