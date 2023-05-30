# Definitions
ROOT                    := $(PWD)
GO_VERSION              := 1.20.4
OS					  	:= linux
ARCH                    := amd64

.PHONY: version
version:
	$(eval VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD))
	@echo "Version: $(VERSION)"

.PHONY: build-docker
build-docker:
	docker build -t mactat/framed -f ./dockerfiles/dev.dockerfile .

.PHONY: build-in-docker
build-in-docker:
	docker run \
		--rm \
		-v $(ROOT):/app \
		--env GOOS=$(OS) \
		--env GOARCH=$(ARCH) \
		golang:$(GO_VERSION)-alpine3.18 \
		/bin/sh -c "cd /app && go build -o ./build/ ./framed.go"
	sudo chown -R $(USER):$(USER) ./build

.PHONY: release-lin
release-lin: version
	$(MAKE) build-in-docker OS=linux ARCH=amd64
	cd ./build && tar -zcvf ./framed-linux-$(VERSION).tar.gz ./framed
	$(MAKE) clean

# make release for windows
.PHONY: release-win
release-win: version
	$(MAKE) build-in-docker OS=windows ARCH=amd64
	cd ./build && tar -zcvf ./framed-windows-$(VERSION).tar.gz ./framed.exe
	$(MAKE) clean

# make release for mac
.PHONY: release-mac
release-mac: version
	$(MAKE) build-in-docker OS=darwin ARCH=amd64
	cd ./build && tar -zcvf ./framed-macos-$(VERSION).tar.gz ./framed
	$(MAKE) clean

.PHONY: build
build:
	go build -o ./build/ ./framed.go

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -f ./build/framed
	rm -f ./build/framed.exe
