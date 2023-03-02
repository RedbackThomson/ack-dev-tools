SHELL := /bin/bash 
GO111MODULE=on

# Build ldflags
VERSION ?= "v0.0.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
PKG_PATH=github.com/aws-controllers-k8s/dev-tools/pkg
GO_LDFLAGS=-ldflags "-X $(PKG_PATH)/version.GitVersion=$(VERSION) \
			-X $(PKG_PATH)/version.GitCommit=$(GITCOMMIT) \
			-X $(PKG_PATH)/version.BuildDate=$(BUILDDATE)"

GO_CMD_FLAGS=-tags codegen -gcflags="all=-N -l" 
GO_CMD_LOCAL_FLAGS=-modfile=go.local.mod $(GO_CMD_FLAGS)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GO_CMD_FLAGS} ${GO_LDFLAGS} -o ./bin/ackdev ./cmd/ackdev/main.go

install: build
	cp ./bin/ackdev $(shell go env GOPATH)/bin/ackdev

test:
	go test -tags $(shell go env GOOS) -v ./...

.PHONY: test install mocks

mocks:
	@echo -n "building mocks for pkg/git ... "
	@mockery --quiet --name OpenCloner --tags=codegen --case=underscore --output=mocks --dir=pkg/git
	@echo "ok."
	@echo -n "building mocks for pkg/github ... "
	@mockery --quiet --all --tags=codegen --case=underscore --output=mocks --dir=pkg/github
	@echo "ok."
