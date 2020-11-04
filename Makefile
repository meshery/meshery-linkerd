GOPATH = $(shell go env GOPATH)

check:
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

.PHONY: local-check
local-check: tidy
local-check: golangci-lint

.PHONY: tidy
tidy:
	@echo "Executing go mod tidy"
	go mod tidy

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)
	@echo
	$(GOPATH)/bin/golangci-lint run

$(GOLANGCILINT):
	(cd /; GO111MODULE=on GOPROXY="direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)

# If your internet was blocked by sth, you can use this proxy
#$(GOLANGCILINT):
#	(cd /; GO111MODULE=on GOPROXY="https://goproxy.cn,direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)
