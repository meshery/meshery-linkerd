IMAGE?=layer5/meshery-linkerd
GOPATH = $(shell go env GOPATH)

GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_STRIPPED_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` | cut -c 2-)

check: error
	golangci-lint run

check-clean-cache:
	golangci-lint cache clean
	
protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-linkerd .

docker-run:
	(docker rm -f meshery-linkerd) || true
	docker run --name meshery-linkerd -d \
	-p 10001:10001 \
	-e DEBUG=true \
	layer5/meshery-linkerd

run:
	DEBUG=true go run main.go

error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

local-check: tidy
local-check: golangci-lint

tidy:
	@echo "Executing go mod tidy"
	go mod tidy

golangci-lint: $(GOLANGCILINT)
	@echo
	$(GOPATH)/bin/golangci-lint run

$(GOLANGCILINT):
	(cd /; GO111MODULE=on GOPROXY="direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

multi: ## Build service image to be deployed as a desktop extension
	docker build --tag=$(IMAGE) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg GIT_STRIPPED_VERSION=$(GIT_STRIPPED_VERSION) .

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'	

.PHONY: error tidy golangci help run local-check