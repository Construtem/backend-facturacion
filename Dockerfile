# Etapa 1: build
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Forzar binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Etapa 2: final
FROM alpine:latest

WORKDIR /root/

# Instalar lib necesaria solo si el binario la requiere
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
