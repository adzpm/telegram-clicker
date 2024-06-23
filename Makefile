PROJECT_NAME = "telegram-clicker"

.PHONY: build
build:
	@echo "Building the project..."
	@GOOS=linux GOARCH=amd64 go build -o bin/$(PROJECT_NAME) cmd/main.go
