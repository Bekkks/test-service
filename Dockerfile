FROM golang:1.21-alpine AS builder

WORKDIR /app

# Установка зависимостей для сборки
RUN apk add --no-cache git

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Генерация Swagger документации
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go || true

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
