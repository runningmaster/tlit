.DEFAULT_GOAL := test

## fmt: formats Go source code
fmt:
	@gofumpt -l -w .
.PHONY: fmt

## vet: examines Go source code and reports suspicious constructs
vet: fmt
	@go vet ./...
	@staticcheck ./...
	@shadow ./...
.PHONY: vet


## lint: runs linters for Go source code
lint: vet
	@golangci-lint --version
	@golangci-lint run ./...
.PHONY: lint

## test: runs all tests
test: lint
	@go list ./... | grep -v vendor/ | xargs -L1 go test -cover -v -count=1
.PHONY: test


## bench: runs `go test` with benchmarks
bench: lint
	@go list ./... | grep -v vendor/ | xargs -L1 go test -bench . -benchmem
.PHONY: bench

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help