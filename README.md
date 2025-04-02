# url-shortener

## Сервис для удобства сокращения ссылок
`go sqlite chi`

### Быстрый старт

```bash
go mod download
CONFIG_PATH=config/local.yaml

go run cmd/url-shortener/main.go
```
*Make*
```bash
make fastRun
```

#chi #sqlite #slog

### Docker
```bash
docker build -t url-shortener .
docker run -it --rm -p 8080:8080 url-shortener
```
*Make*
```bash
make buildAndRun
```

### HRM
```bash
go install github.com/air-verse/air@latest
air
```

*Make*
```bash
make airStart
```

### Swagger
Запустить проект
```bash
make fastRun
```
Открыть по [ссылке](http://localhost:8080/swagger/index.html#/url/post_url)

После внесения правок:
```bash
make initSwaggerDoc
```