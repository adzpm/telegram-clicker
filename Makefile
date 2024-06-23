.PHONY: build
build:
	@echo "Building the project..."
	GOOS=linux GOARCH=amd64 go build -o bin/tgc.bin cmd/main.go
