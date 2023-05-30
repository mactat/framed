# Definitions
ROOT                    := $(PWD)
GO_VERSION              := 1.20.4

.PHONY: build-docker
build-docker:
	docker build -t mactat/framed -f ./dockerfiles/dev.dockerfile .

.PHONY: release-lin
release-lin:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && go build -o ./build/framed ./main.go"
	sudo chown -R $(USER):$(USER) ./build

# make release for windows
.PHONY: release-win
release-win:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && GOOS=windows GOARCH=amd64 go build -o ./build/framed.exe ./main.go"
	sudo chown -R $(USER):$(USER) ./build

# make release for mac
.PHONY: release-mac
release-mac:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && GOOS=darwin GOARCH=amd64 go build -o ./build/framed ./main.go"
	sudo chown -R $(USER):$(USER) ./build

.PHONY: build
build:
	go build -o ./build/framed ./main.go

.PHONY: format
format:
	go fmt ./...
