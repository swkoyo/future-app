all: build

build:
	@echo "Building..."
	@go build -o tmp/main cmd/api/main.go

run:
	@go run cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -rf tmp

watch:
	air

.PHONY: all build run clean watch
