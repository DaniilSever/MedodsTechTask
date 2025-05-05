FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /src
RUN adduser --disabled-password --gecos "" appuser && chown -R appuser /src

# Копируем vendor-архив
COPY src/vendor.tar.gz .
RUN tar -xzf vendor.tar.gz && rm vendor.tar.gz

COPY src/go.mod src/go.sum ./

# Copyng source code
COPY src/ /src/

# Сборка с локальными зависимостями
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o /app/server ./app/main.go

FROM cgr.dev/chainguard/wolfi-base:latest

WORKDIR /app

COPY --from=builder /app/server /server

EXPOSE 8080
CMD ["/server"]