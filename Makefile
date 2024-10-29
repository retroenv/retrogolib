GOLANGCI_VERSION = v1.61.0

help: ## show help, shown by default if no target is specified
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## run code linters
	golangci-lint run

build-all: ## build code with all 3 GUI mode settings
	go build ./...
	go build -tags noopengl,sdl ./...
	go build -tags nogui ./...

test: ## run tests
	go test -timeout 10s -race ./...

test-no-gui: ## run unit tests with gui disabled
	go test -timeout 10s -tags nogui ./...

test-coverage: ## run unit tests and create test coverage
	go test -timeout 10s -tags nogui ./... -coverprofile .testCoverage -covermode=atomic -coverpkg=./...

test-coverage-web: test-coverage ## run unit tests and show test coverage in browser
	go tool cover -func .testCoverage | grep total | awk '{print "Total coverage: "$$3}'
	go tool cover -html=.testCoverage

install-linters: ## install all used linters
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION}
