OAPI_CODEGEN_OUTPUT := internal/adapters/api/gen_api.go

generate-api:
	go tool oapi-codegen -generate types,server -package api openapi.yaml > $(OAPI_CODEGEN_OUTPUT)
	@echo "API code generated at $(OAPI_CODEGEN_OUTPUT)."

run:
	docker-compose up --build -d
	@echo "Service is running at "

stop:
	docker-compose down

logs:
	docker-compose logs -f

logs-app:
	docker-compose logs -f app

logs-db:
	docker-compose logs -f db

install-tools:
	@echo "Installing dev-tools (oapi-codegen, golangci-lint)..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-deps: install-tools
	go mod tidy
	go mod verify


fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

help:
	@echo "All commands:"
	@echo "    make run           - start in Docker"
	@echo "    make stop          - stop docker process"
	@echo "    make logs          - show all logs"
	@echo "    make logs-app      - show application logs"
	@echo "    make logs-db       - show database logs"
	@echo "    make install-tools - install development tools"
	@echo "    make install-deps  - install dependencies (dev-tools included)"
	@echo "    make generate-api  - generate api code according to openapi scheme"
	@echo "    make fmt           - format all code"
	@echo "    make lint          - lint code"
