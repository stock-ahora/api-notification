# Build stage
FROM golang:1.25-alpine AS builder

# Establecer directorio de trabajo
WORKDIR /build

# Copiar archivos de dependencias
COPY go.mod go.sum* ./

# Descargar dependencias y actualizar paquetes del sistema
RUN apk update && apk upgrade --no-cache && go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./

# Final stage - Imagen mínima
FROM scratch

# Copiar el binario compilado
COPY --from=builder /build/api /api

# Exponer puerto
EXPOSE 8084

# Ejecutar el binario
ENTRYPOINT ["/api"]