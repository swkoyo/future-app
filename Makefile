all: build

build:
	@echo "Building..."
	@go build -o tmp/main cmd/api/main.go

run:
	@go run cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -rf tmp

seed:
	@echo "Seeding..."
	@go run cmd/scripts/seed/main.go seed

watch:
	air

.PHONY: all build run clean watch seed
