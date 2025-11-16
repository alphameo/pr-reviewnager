OAPI_CODEGEN_OUTPUT := internal/adapters/api/gen_api.go

generate-api:
	go tool oapi-codegen -generate types,server -package api openapi.yaml > $(OAPI_CODEGEN_OUTPUT)
