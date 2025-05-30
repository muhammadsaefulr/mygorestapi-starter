include .env
export $(shell sed 's/=.*//' .env)

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
	@swag init -g cmd/main.go

migration-%:
	@migrate create -ext sql -dir pkg/infrastructure/persistence/migrations create-table-$(subst :,_,$*)

migrate-up:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path pkg/infrastructure/persistence/migrations up

migrate-down:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path pkg/infrastructure/persistence/migrations down

migrate-docker-up:
	@docker run -v ./pkg/infrastructure/persistence/migrations:/migrations --network NimeStreamAPI_go-network migrate/migrate -path=/migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-docker-down:
	@docker run -v ./pkg/infrastructure/persistence/migrations:/migrations --network NimeStreamAPI_go-network migrate/migrate -path=/migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

docker:
	@chmod -R 755 ./database/init
	@docker-compose up --build

docker-test:
	@docker-compose up -d && make tests

docker-down:
	@docker-compose down --rmi all --volumes --remove-orphans

docker-cache:
	@docker builder prune -f
