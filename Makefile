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

