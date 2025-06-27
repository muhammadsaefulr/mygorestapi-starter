include .env
export $(shell sed 's/=.*//' .env)

generator:
ifndef gen
	@echo "Usage: make generator gen=nama_modul"
	@exit 1
else
	@echo "Generating module: $(gen)"
	@go run cmd/generator/generate_base_service.go generate module $(gen)
endif

rem-generator:
ifndef gen
	@echo "Usage: make rem-generator gen=moduleName"
	@exit 1
else
	@echo "Undoing generated module: $(gen)"
	@rm -fv \
		internal/delivery/http/controller/$(gen)_controller/$(gen)_controller.go \
		internal/delivery/http/router/$(gen)_router.go \
		internal/repository/$(gen)/$(gen)_repository.go \
		internal/repository/$(gen)/$(gen)_repository_impl.go \
		internal/service/$(gen)_service/$(gen)_service.go \
		internal/service/$(gen)_service/$(gen)_service_impl.go \
		internal/shared/convert_types/$(gen)_converter.go \
		internal/domain/dto/$(gen)/request/request.go \
		internal/domain/dto/$(gen)/response/response.go \
		internal/domain/model/$(gen)_model.go
	@echo "Checking for empty folders to remove..."
	@find internal -type d -empty -delete
endif

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
