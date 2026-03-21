## Integration Test Makefile for Proxmox Backup Server Terraform Provider

.PHONY: test test-integration test-unit build clean help setup-dev

# Default target
help:
	@echo "Available targets:"
	@echo "  build             - Build the terraform provider"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests (requires PBS_* env vars)"
	@echo "  clean             - Clean build artifacts"
	@echo "  setup-dev         - Setup development environment"
	@echo ""
	@echo "Integration test environment variables:"
	@echo "  PBS_ENDPOINT      - PBS server endpoint (e.g., https://pbs.example.com:8007)"
	@echo "  PBS_USERNAME      - PBS username"
	@echo "  PBS_PASSWORD      - PBS password"
	@echo "  PBS_INSECURE      - Set to true for self-signed certificates"
	@echo ""
	@echo "Example integration test run:"
	@echo "  PBS_ENDPOINT=https://pbs.local:8007 PBS_USERNAME=admin PBS_PASSWORD=secret PBS_INSECURE=true make test-integration"

# Build the provider
build:
	go build -o terraform-provider-pbs .

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	go test -v -short ./...

# Run integration tests
test-integration:
	../testacc

# Clean build artifacts
clean:
	rm -f terraform-provider-pbs
	rm -rf test/terraform.tfstate*
	rm -rf test/.terraform/

# Setup development environment
setup-dev:
	go mod tidy
	go mod download
