SHELL:=/bin/bash
LAMBDA_BIN_NAME=main
ZIP_FILENAME=catalyst-main.zip
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_FILES=$(shell go list ./... | grep -v /vendor/)

.SILENT:

all: build zip
	echo -e "\e[32mUpload $(ZIP_FILENAME) to your AWS Lambda.\nThe handler is \"$(LAMBDA_BIN_NAME)\".\e[0m"; \
	rm -f $(LAMBDA_BIN_NAME) 2> /dev/null

build:
	echo -e "\e[32mBuild Go binary ($(LAMBDA_BIN_NAME))\e[0m"; \
	$(GO_BUILD_ENV) go build -v -o $(LAMBDA_BIN_NAME) .

zip:
	zip $(ZIP_FILENAME) $(LAMBDA_BIN_NAME)
