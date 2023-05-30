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
	$(eval VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD))
	cd ./build && tar -zcvf ./framed-linux-$(VERSION).tar.gz ./framed
	rm ./build/framed

# make release for windows
.PHONY: release-win
release-win:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && GOOS=windows GOARCH=amd64 go build -o ./build/framed.exe ./main.go"
	sudo chown -R $(USER):$(USER) ./build
	$(eval VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD))
	cd ./build && tar -zcvf ./build/framed-win-$(VERSION).tar.gz ./build/framed.exe
	rm -rf ./build/framed.exe

# make release for mac
.PHONY: release-mac
release-mac:
	docker run --rm -v $(ROOT):/app golang:$(GO_VERSION)-alpine3.18 /bin/sh -c "cd /app && GOOS=darwin GOARCH=amd64 go build -o ./build/framed ./main.go"
	sudo chown -R $(USER):$(USER) ./build
	$(eval VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD))
	cd ./build && tar -zcvf ./build/framed-mac-$(VERSION).tar.gz ./build/framed
	rm ./build/framed

.PHONY: build
build:
	go build -o ./build/framed ./main.go

.PHONY: format
format:
	go fmt ./...
