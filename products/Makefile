.PHONY: gen-docs
gen-docs:
	@swag init -g ./cmd/api/main.go -d cmd,internal && swag fmt

.PHONY: build
build:
	@go build -v ./...

.PHONY: test
test:
	@go test -v ./...
