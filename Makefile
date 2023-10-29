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

lint: vet
	golint $(http)
	golint $(db)
.PHONY:vet

run: lint
	go run $(http)
.PHONY:run

build:
	go build $(http)
.PHONY:build
