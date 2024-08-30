# To run unit tests
tests:
	@echo "Running tests"
	go test -count=1 -race ./...

# To generate static files
generate-doc:
	@mkdir -p ./main/static
	go run ./main/main.go -generate

.PHONY: tests generate-static