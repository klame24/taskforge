.PHONY: run build docker-build docker-up deps lint test

deps:
	go mod tidy

run: deps
	go run ./cmd/api

build: deps
	go build -o bin/api ./cmd/api

docker-build: deps
	docker build -f docker/Dockerfile -t taskforge-api .

docker-up: deps
	docker compose -f docker/docker-compose.dev.yml up --build

lint: deps
	golangci-lint run

test:
	go test ./... -v

clean:
	rm -rf bin/