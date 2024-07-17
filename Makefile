all: build

build:
	@echo "Building..."
	@go build -o tmp/main cmd/api/main.go

run:
	@go run cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -rf tmp

test:
	@echo "Testing..."
	@go test -v ./...

seed:
	@echo "Seeding..."
	@go run cmd/scripts/seed/main.go seed

watch:
	air

.PHONY: build run clean test seed watch
