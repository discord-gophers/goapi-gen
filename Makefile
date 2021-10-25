help:
	@echo "This is a helper makefile for goapi-gen"
	@echo "Targets:"
	@echo "    generate:    regenerate all generated files"
	@echo "    test:        run all tests"

generate:
	go generate ./pkg/...
	go generate ./...

test:
	go test -cover ./...
