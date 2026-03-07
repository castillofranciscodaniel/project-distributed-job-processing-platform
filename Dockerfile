# --- STAGE 1: Constructor ---
FROM golang:1.24-alpine AS builder

# Instalamos dependencias necesarias para la compilación
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Aprovechar cache de Docker para dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente
COPY . .

# Compilamos la API y el Worker por separado
RUN CGO_ENABLED=0 GOOS=linux go build -o api-bin ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-bin ./cmd/worker/main.go

# --- STAGE 2: Ejecución (Imagen final ligera) ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Traemos los binarios compilados
COPY --from=builder /app/api-bin .
COPY --from=builder /app/worker-bin .

# Por seguridad: Correr como usuario no-root
RUN adduser -D appuser && chown -R appuser /app
USER appuser

# Exponemos puerto del API
EXPOSE 8080

# Por defecto arranca la API
# En ECS, para el Worker, sobreescribiremos el CMD a ["./worker-bin"]
CMD ["./api-bin"]