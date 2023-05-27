# Definitions
ROOT                    := $(PWD)
GO_VERSION              := 1.20.4

.PHONY: build-docker
build-docker:
	docker build -t mactat/framed -f ./dockerfiles/dev.dockerfile .

.PHONY: release
release:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && go build -o ./build/framed ./main.go"

.PHONY: build
build:
	go build -o ./build/framed ./main.go

.PHONY: format
format:
	go fmt ./...
