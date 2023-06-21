# Definitions
ROOT                    := $(PWD)
GO_VERSION              := 1.20.4
ALPINE_VERSION          := 3.18
OS					  	:= linux
ARCH                    := amd64
BUILD                   := true

.PHONY: version
version:
	$(eval VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD))
	@echo "Version: $(VERSION)"

.PHONY: build-docker
build-docker:
	docker build \
	-t mactat/framed \
	-f ./dockerfiles/dockerfile \
	--build-arg GO_VERSION=$(GO_VERSION) \
	--build-arg ALPINE_VERSION=$(ALPINE_VERSION) \
	.

.PHONY: release-docker
release-docker: version build-docker
	docker tag mactat/framed mactat/framed:$(VERSION)
	docker tag mactat/framed mactat/framed:alpine-$(ALPINE_VERSION)-$(VERSION)
	docker push mactat/framed:$(VERSION)
	docker push mactat/framed:alpine-$(ALPINE_VERSION)-$(VERSION)
	docker push mactat/framed:latest

.PHONY: build-in-docker
build-in-docker: clean version
	docker run \
		--rm \
		-v $(ROOT):/app \
		--env GOOS=$(OS) \
		--env GOARCH=$(ARCH) \
		golang:$(GO_VERSION)-alpine$(ALPINE_VERSION) \
		/bin/sh -c "cd /app && go build -o ./build/ ./framed.go"
	sudo chown -R $(USER):$(USER) ./build
	cd ./build && \
	tar -zcvf \
	  ./framed-$(OS)-$(ARCH)-$(VERSION).tar.gz \
	  ./framed$(if $(filter $(OS),windows),.exe)

.PHONY: release-lin
release-lin:
	$(MAKE) build-in-docker OS=linux ARCH=amd64

.PHONY: release-win
release-win:
	$(MAKE) build-in-docker OS=windows ARCH=amd64

.PHONY: release-mac
release-mac:
	$(MAKE) build-in-docker OS=darwin ARCH=amd64

.PHONY: build
build:
	go build -o ./build/ ./framed.go

.PHONY: test
test:
	@if [ "$(BUILD)" = "true" ]; then\
        $(MAKE) build-docker;\
    fi
	docker build -f ./dockerfiles/test.dockerfile -t framed-test .
	@if [ "$(EXPORT)" = "true" ]; then\
		mkdir -p ./results;\
		docker run --rm framed-test /bin/sh -c "/test/bats/bin/bats -F junit /test/" > ./results/test.xml;\
	else\
		docker run --rm framed-test /bin/sh -c "/test/bats/bin/bats --pretty /test/";\
	fi
.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -f ./build/framed
	rm -f ./build/framed.exe
