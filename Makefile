.PHONY: build run clean test

# Build the game
build:
	go build -o bin/ascii-type ./cmd/game

# Run the game
run: build
	./bin/ascii-type

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/ascii-type-linux-amd64 ./cmd/game
	GOOS=darwin GOARCH=amd64 go build -o bin/ascii-type-darwin-amd64 ./cmd/game
	GOOS=windows GOARCH=amd64 go build -o bin/ascii-type-windows-amd64.exe ./cmd/game
