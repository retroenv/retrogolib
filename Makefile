GOLANGCI_VERSION = v2.13.2
RETROGOLINT_VERSION = v1.0.3
TEST_TIMEOUT ?= 60s

help: ## show help, shown by default if no target is specified
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## run code linters
	golangci-lint run
	retrogolint

build: ## build code
	CGO_ENABLED=0 go build ./...

test: ## run tests
	go test -short -timeout $(TEST_TIMEOUT) -race ./...

test-6502-release: ## qualify NMOS legal opcodes and functional binaries using required pinned corpora
	CPU6502_QUALIFY=1 go test -race -count=1 -v -timeout 10m -run '^TestRelease(NMOS|Dormann)$$' ./arch/cpu/cpu6502/

test-integration: ## run long-running CPU integration tests
	go test -v -run 'TestSingleStep|TestDormann' -timeout 0 -race ./arch/cpu/cpu6502/
	go test -tags singlestep -v -run TestSingleStep -timeout 0 -race ./arch/cpu/cpu65816/
	go test -tags singlestep -v -run TestSingleStep -timeout 0 -race ./arch/cpu/cpu68000/
	go test -v -run TestSingleStep -timeout 0 -race ./arch/cpu/sm83/
	go test -v -run 'TestSingleStep|TestZexdoc|TestZexall' -timeout 0 -race ./arch/cpu/z80/

test-coverage: ## run unit tests and create test coverage
	go test -short -timeout $(TEST_TIMEOUT) ./... -coverprofile coverage.txt

test-coverage-web: test-coverage ## run unit tests and show test coverage in browser
	go tool cover -func coverage.txt | grep total | awk '{print "Total coverage: "$$3}'
	go tool cover -html=coverage.txt

install-linters: ## install all used linters
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_VERSION}
	go install github.com/retroenv/retrogolint/cmd/retrogolint@${RETROGOLINT_VERSION}
