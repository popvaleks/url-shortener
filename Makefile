setEnv:
	CONFIG_PATH=config/local.yaml

updateDep:
	go mod tidy

downloadDep:
	go mod download

run:
	go run cmd/url-shortener/main.go


buildDocker:
	docker build -t url-shortener .

runDocker:
	docker run -p 8080:8080 url-shortener

runHmr:
	air

initSwaggerDoc:
	swag init -g cmd/url-shortener/main.go \
      --exclude ./config,./docs,./tmp,./storage \
      --parseDependency \
      --parseInternal \
      --parseDepth 2

fastRun:
	make downloadDep
	make run

buildAndRun:
	make buildDocker
	make runDocker

airStart:
	go install github.com/air-verse/air@latest
	air