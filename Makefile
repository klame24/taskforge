run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

docker-build:
	docker build -f docker/Dockerfile -t taskforge-api .

docker-up:
	docker compose -f docker/docker-compose.dev.yml up --build
