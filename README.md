# url-shortener

## Сервис для удобства сокращения ссылок

### Быстрый старт

```bash
go mod download
CONFIG_PATH=config/local.yaml

go run cmd/url-shortener/main.go
```

#chi #sqlite #slog

### Docker
для macOs ARM
```bash
docker build --platform linux/arm64 -t url-shortener .
docker run -it --rm -p 8080:8080 url-shortener
```
или с указанием платформы
```bash
docker buildx build --platform ВАША_ПЛАТФОРМА -t url-shortener .
```
