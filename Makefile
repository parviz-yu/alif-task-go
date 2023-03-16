.PHONY: build
build:docker migrate
	go build -v ./cmd/wallet

.PHONY: docker
docker:
	docker run --name=alif-postgrs -e POSTGRES_PASSWORD='alif' -p 5432:5432 -d --rm postgres
	docker run --name=alif-redis -p 6379:6379 -d --rm redis
	sleep 5

.PHONY: migrate
migrate:
	migrate -path ./schema -database 'postgres://postgres:alif@localhost:5432/postgres?sslmode=disable' up	

.DEFAULT_GOAL := build