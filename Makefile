.PHONY: generate-openapi


generate-openapi: generate-swagger rename-swagger-json rename-swagger-yml

generate-swagger:
	swag init -g ./cmd/server/main.go --output docs --parseDependency --exclude internal/storage/*

rename-swagger-json:
	mv docs/swagger.json docs/kata-peya-openapi.json

rename-swagger-yml:
	mv docs/swagger.yaml docs/kata-peya-openapi.yml

test.integration:
	go test -v ./internal/test/integration/... -count=1 -tags=integration

test.e2e:
	go test -v ./internal/test/e2e/... -count=1 -tags=e2e

test.unit:
	go test -v ./... -count=1