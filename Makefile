include .env
export $(shell sed 's/=.*//' .env)

generator:
	@echo "Usage: make generator NAME=yourModule"
	@exit 1

generator-%:
	@echo "Generating base service for: $(subst :,_,$*)"
	@go run cmd/generator/generate_base_service.go generate module $(subst :,_,$*)

start:
	@go run cmd/main.go

lint:
	@golangci-lint run

tests:
	@go test -v ./test/...

tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')

testsum:
	@cd test && gotestsum --format testname

swag:
	@swag init -g cmd/main.go --parseDependency --parseInternal

migration-%:
	@migrate create -ext sql -dir internal/infrastructure/persistence/migrations create-table-$(subst :,_,$*)

migrate-up:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/infrastructure/persistence/migrations up

migrate-down:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/infrastructure/persistence/migrations down

migrate-docker-up:
	@docker run -v ./internal/infrastructure/persistence/migrations:/migrations --network NimeStreamAPI_go-network migrate/migrate -path=/migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-docker-down:
	@docker run -v ./internal/infrastructure/persistence/migrations:/migrations --network NimeStreamAPI_go-network migrate/migrate -path=/migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

docker:
	@chmod -R 755 ./database/init
	@docker-compose up --build

docker-test:
	@docker-compose up -d && make tests

docker-down:
	@docker-compose down --rmi all --volumes --remove-orphans

docker-cache:
	@docker builder prune -f
