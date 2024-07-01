.PHONY: build
build:
	@echo "building the binary"
	GOOS=linux GOARCH=amd64 go build -o bin/app.bin cmd/main.go
