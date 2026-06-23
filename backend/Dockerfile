# Etapa de construcción
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos del módulo y descargar dependencias
COPY go.mod ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar la aplicación de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa de ejecución
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario compilado desde el builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
