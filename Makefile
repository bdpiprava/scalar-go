# To run unit tests
ROOT_DIR := $(shell pwd)

tests:
	@echo "Running tests"
	go test -count=1 -race ./...

# To generate static files
generate-doc:
	@echo "Generating docs"
	@mkdir -p $(ROOT_DIR)/build/docs
	go run $(ROOT_DIR)/docs/main.go -generate

.PHONY: tests generate-static