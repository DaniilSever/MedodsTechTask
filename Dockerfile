FROM cgr.dev/chainguard/go:latest AS builder

# Устанавливаем swag и validator
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/go-playground/validator/v10@latest

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger
RUN swag init -g src/app/main.go

# Билдим приложение (статически линкуем)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/bin/server ./src/app/main.go

# === Финальная стадия ===
FROM cgr.dev/chainguard/wolfi-base:latest

WORKDIR /app

# Копируем бинарник и документацию
COPY --from=builder /app/bin/server .
COPY --from=builder /app/docs ./docs

# Настраиваем пользователя (не root)
USER 65532:65532

EXPOSE 8080
CMD ["./server"]