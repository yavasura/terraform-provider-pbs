NAME=terraform-provider-pbs
VERSION ?= $(shell tr -d '\n' < VERSION)
TEST_PROVIDER_VERSION ?= ${VERSION}
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/registry.terraform.io/yavasura/pbs/${VERSION}/${GOOS}_${GOARCH}

default: build

.PHONY: build
build:
	go build -o ${NAME} .

.PHONY: install
install: build
	mkdir -p ${PLUGIN_DIR}
	cp ${NAME} ${PLUGIN_DIR}/

.PHONY: test
test:
	go test ./...

.PHONY: test-unit
test-unit:
	go test ./...

.PHONY: testacc
testacc:
	./testacc

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	gofmt -s -w .
	go mod tidy

.PHONY: clean
clean:
	rm -f ${NAME}

.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: release-test
release-test:
	goreleaser release --snapshot --clean --skip=publish

.PHONY: release
release:
	@echo "To create a release:"
	@echo "  1. Ensure all changes are committed"
	@echo "  2. Create and push a tag: git tag -a v${VERSION} -m 'Release v${VERSION}'"
	@echo "  3. Push the tag: git push origin v${VERSION}"
	@echo "  4. GitHub Action will automatically build and publish"


.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the provider binary"
	@echo "  install       - Build and install the provider locally"
	@echo "                  Variables: VERSION, GOOS, GOARCH"
	@echo "  test          - Run unit tests"
	@echo "  testacc       - Run acceptance tests"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Remove built binaries"
	@echo "  docs          - Generate documentation"
	@echo "  release-test  - Test release build locally (requires goreleaser)"
	@echo "  release       - Show release instructions"
	@echo "  help          - Show this help message"
