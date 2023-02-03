.DEFAULT_GOAL := test

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go vet ./...
.PHONY:vet

test: vet
	go test ./...
.PHONY:test