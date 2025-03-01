SWAG ?= swag

.PHONY: swagger-docs-packer test

swagger-docs-packer:
	@echo "Generating Swagger documentation for Packer..."
	$(SWAG) init -d ./cmd/packer,./pkg/server --parseDependency --output ./docs/packer

test:
	go test -v ./...

run:
	go run cmd/packer/main.go
