.DEFAULT_GOAL := run

http := ./cmd/http
db := ./album

fmt:
	go fmt $(http)
	go fmt $(db)
.PHONY:fmt

vet: fmt
	go vet $(http)
	go vet $(db)
.PHONY:vet

run: vet
	go run $(http)
.PHONY:run

build: 
	go build $(http)
.PHONY:build
