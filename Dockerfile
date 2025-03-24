# Аргумент для выбора платформы (по умолчанию linux/arm64 для M1/M2/M3)
ARG TARGETPLATFORM=linux/arm64

# Этап сборки
FROM --platform=$TARGETPLATFORM golang:1.23.3 AS builder

# Устанавливаем зависимости для SQLite
RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Определяем архитектуру для сборки
RUN case ${TARGETPLATFORM} in \
    "linux/amd64")  GOARCH="amd64"  ;; \
    "linux/arm64")  GOARCH="arm64"  ;; \
    *)              GOARCH="arm64"  ;; \
    esac && \
    CGO_ENABLED=1 GOARCH=$GOARCH go build -ldflags="-s -w" -o url-shortener ./cmd/url-shortener

# Этап запуска
FROM --platform=$TARGETPLATFORM alpine:latest

# Устанавливаем runtime-зависимости для SQLite
RUN apk add --no-cache sqlite-dev && \
    mkdir -p /app/storage && \
    touch /app/storage/storage.db && \
    chmod 666 /app/storage/storage.db

WORKDIR /app

# Копируем бинарник и конфиги
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/config/docker.yaml /app/config/
# Внимание !!! Копируется БД
COPY --from=builder /app/storage/storage.db /app/storage/

ENV CONFIG_PATH=/app/config/docker.yaml

EXPOSE 8080
CMD ["./url-shortener"]