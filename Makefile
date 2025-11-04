# To run unit tests
ROOT_DIR := $(shell pwd)
# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

tests:
	@echo "$(OK_COLOR)==> Running tests...$(NO_COLOR)"
	@go install gotest.tools/gotestsum@latest
	@gotestsum --format=testname -- -v -race=1 -coverprofile=coverage_unit.txt -coverpkg=./...

# To generate static files
generate-doc:
	@echo "$(OK_COLOR)==> Generating docs...$(NO_COLOR)"
	@mkdir -p $(ROOT_DIR)/build/docs
	go run $(ROOT_DIR)/docs/main.go -generate

.PHONY: tests generate-static