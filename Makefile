test:
	@echo "Testing Go packages..."
	@go test ./... -cover

test-short:
	@echo "Testing Go packages..."
	@go test ./... -cover -short

mocks:
	@echo "Regenerate mocks..."
	@go generate ./...