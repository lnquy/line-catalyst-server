SHELL:=/bin/bash
PROJECT_NAME=line-catalyst-server
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_FILES=$(shell go list ./... | grep -v /vendor/)

BUILD_VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TAG=$(BUILD_VERSION)-$(GIT_COMMIT)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_TAG)

.SILENT:

all: fmt vet install test

build:
	echo -e "\e[32mBuild Go binary ($(PROJECT_NAME)-$(BUILD_TAG).bin)\e[0m"; \
	git --no-pager log -1; \
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_TAG).bin .

install:
	$(GO_BUILD_ENV) go install

vet:
	go vet $(GO_FILES)

lint:
	golint $(GO_FILES)

fmt:
	go fmt $(GO_FILES)

test:
	go test $(GO_FILES) -cover

vendor:
	govendor sync

integration_test:
	go test -tags=integration $(GO_FILES)

#compose: build
#	mv $(PROJECT_NAME)-$(BUILD_TAG).bin $(PROJECT_NAME).bin; \
#	cd deployment/docker && docker-compose up
#
#docker: build
#	echo -e "\n\e[32mBuild Docker image ($(DOCKER_IMAGE))\e[0m"; \
#	cd deployment/docker; \
#	mv ../../$(PROJECT_NAME)-$(BUILD_TAG).bin $(PROJECT_NAME).bin; \
#	docker build -t $(DOCKER_IMAGE) .; \
#	rm -f $(PROJECT_NAME).bin 2> /dev/null; \
#	cd ../..
