FROM golang:1.22-alpine AS builder

# Устанавливаем swag и зависимости
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/go-playground/validator/v10@latest

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY src/go.mod src/go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger-документацию
RUN swag init -g src/app/main.go  # Укажи путь к main.go

# Билдим приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./src/app/main.go

# === Финальная стадия (минимальный образ) ===
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и документацию
COPY --from=builder /app/bin/server .
COPY --from=builder /app/docs ./docs

# Порт для приложения (Gin/Echo)
EXPOSE 8080

# Запускаем сервер
CMD ["./server"]