# Imagen base de Go
FROM golang:1.23-alpine

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Construir el binario
RUN go build -o server .

# Exponer el puerto en el que correrá el servidor
EXPOSE 5080

# Ejecutar el servidor
CMD ["./server"]
