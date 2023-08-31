.DEFAULT_GOAL := run

fmt:
	go fmt ./cmd/http
.PHONY:fmt

vet: fmt
	go vet ./cmd/http
.PHONY:vet

run: vet
	go run ./cmd/http
.PHONY:run

build: 
	go build ./cmd/http
.PHONY:build
