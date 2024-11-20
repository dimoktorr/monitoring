# Set shell
SHELL=/bin/bash

# Variables
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
PROTOC_OS := $(OS)
ifeq ($(PROTOC_OS), darwin)
	PROTOC_OS := "osx"
endif
ARCH := $(shell uname -m)
PROTOC_ARCH := $(ARCH)
ifeq ($(PROTOC_ARCH), arm64)
	PROTOC_ARCH := "aarch_64"
endif

PROTOC_VERSION := 28.0

proto-gen: protobuf-install
	protoc -I . \
		-I ./api \
		--go_out ./pkg/ --go_opt paths=source_relative \
		--go-grpc_out ./pkg/ --go-grpc_opt paths=source_relative \
		api/v1/*.proto

protobuf-install: require
ifeq ($(shell ${GOPATH}/bin/protoc --version 2>/dev/null), )
	@echo "Installing protoc for ${PROTOC_OS} ${PROTOC_ARCH}"
	@curl -LO https://nexus.services.mts.ru/repository/github-raw-proxy/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip
	@unzip protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip -d ${GOPATH}/
	@rm -f ${GOPATH}/readme.txt && rm -f protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
endif

require:
	@mkdir -p ./bin
	@which go &> /dev/null || (echo "error: go is required."; exit 1)
	@which curl &> /dev/null || (echo "error: curl is required."; exit 1)
	@which git &> /dev/null || (echo "error: git is required."; exit 1)

protobuf-uninstall: require
ifneq ($(shell ${GOPATH}/bin/protoc --version 2>/dev/null), )
	@rm -f ${GOPATH}/bin/protoc && rm -f ${GOPATH}/bin/protoc-gen* && rm -rf ${GOPATH}/include/google/protobuf
	@echo "protoc uninstalled"
endif