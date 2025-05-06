FROM golang:1.24 AS builder

ENV http_proxy http://213.176.75.31:3128
ENV https_proxy http://213.176.75.31:3128

ARG ENV=${ENV}
ENV ENV=${ENV}

WORKDIR /src

# Копируем go.mod и go.sum
COPY src/go.mod src/go.sum ./

# Установка зависимостей
RUN go mod tidy

RUN go get github.com/jackc/pgx/v5
RUN go get github.com/jackc/pgx/v5/pgxpool
RUN go get github.com/jackc/pgx/v5/pgconn
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/cors
RUN go get github.com/swaggo/files
RUN go get github.com/swaggo/gin-swagger
RUN go get github.com/swaggo/swag@latest


# Установка зависимостей
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем исходники
COPY ./src /src

# Стадия запуска
FROM golang:1.24 AS runner

ENV http_proxy=http://213.176.75.31:3128
ENV https_proxy=http://213.176.75.31:3128

WORKDIR /src

# Копируем проект и инструменты
COPY --from=builder /src /src
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/swag /usr/local/bin/swag

ENV ENV=local
ENV PATH="/go/bin:$PATH"

EXPOSE 8080
CMD ["air"]