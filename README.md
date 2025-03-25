# url-shortener

## Сервис для удобства сокращения ссылок
`go sqlite chi`

### Быстрый старт

```bash
go mod download
CONFIG_PATH=config/local.yaml

go run cmd/url-shortener/main.go
```

#chi #sqlite #slog

### Docker
```bash
docker build -t url-shortener .
docker run -it --rm -p 8080:8080 url-shortener
```

### HRM
```bash
go install github.com/air-verse/air@latest
air
```