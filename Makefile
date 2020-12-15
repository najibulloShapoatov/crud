
build:
	go build ./cmd


run:
	go run ./cmd/main.go


.PHONY: run build

.DEFAULT_GOAL := run