## TODO Need to be fix Does not working well on ubuntu &_&
########################################
GOPATH        = $(SHELL go env GOPATH)
GOLINT        = /root/go/bin/golint
GOIMPORTS     = /root/go/bin/goimports
# MISSPELL      = $(GOPATH)/bin/misspell
GOSEC         = /root/go/bin/gosec
ERRCHECK      = /root/go/bin/errcheck
STATICCHECK   = /root/go/bin/staticcheck
GOCYCLO       = /root/go/bin/gocyclo
ARCH          = $(SHELL uname -p)
########################################

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
GOFLAGS    :=

SHELL      = /usr/bin/env bash


GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)


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

local-check:
	@echo
	go mod tidy
	@echo
	go mod verify


.PHONY: check
check: lint
# check: misspell
check: gosec
check: err-check
# Comment out the test, need to update the Meshery's protobuf's configurations
# check: static-check
check: vet
check: gocyclo

.PHONY: lint
lint: $(GOLINT)
	@echo
	@echo "==> Check & Review code <=="
	GO111MODULE=auto GOPROXY=direct go list ./... | grep -v /vendor/ | xargs -L1 $(GOLINT) -set_exit_status

# .PHONY: misspell
# misspell:
# 	@echo
# 	@echo "==> Correct commonly misspelled English words <=="
# 	GO111MODULE=auto GOPROXY=direct $(MISSPELL) transform pkg

# https://github.com/securego/gosec/issues/458
.PHONY: gosec
gosec: $(GOSEC)
	@echo
	@echo "==> Inspects source code for security problems <=="
	GO111MODULE=auto GOPROXY=direct $(GOSEC) ./...

.PHONY: err-check
err-check: $(ERRCHECK)
	@echo
	@echo "==> Error check <=="
	GO111MODULE=auto GOPROXY=direct $(ERRCHECK) ./...

.PHONY: static-check
static-check: $(STATICCHECK)
	@echo
	@echo "==> Static check <=="
	GO111MODULE=auto GOPROXY=direct $(STATICCHECK) -checks all,-ST1000 ./...

.PHONY: vet
vet:
	@echo
	@echo "==> Vet <=="
	GO111MODULE=auto GOPROXY=direct go vet ./...

# https://github.com/fzipp/gocyclo
.PHONY: gocyclo
gocyclo: $(GOCYCLO)
	@echo
	@echo "==> GOCYCLO <=="
	GO111MODULE=auto GOPROXY=direct $(GOCYCLO) .


$(ERRCHECK):
	(cd /; GO111MODULE=auto GOPROXY=direct go get -u github.com/kisielk/errcheck)

$(STATICCHECK):
	(cd /; GO111MODULE=auto GOPROXY=direct go get -u honnef.co/go/tools/cmd/staticcheck)

$(GOCYCLO):
	(cd /; GO111MODULE=auto GOPROXY=direct go get -u github.com/fzipp/gocyclo)

$(GOSEC):
	(cd /; GO111MODULE=auto GOPROXY=direct go get github.com/securego/gosec/cmd/gosec)

$(GOLINT):
	(cd /; GO111MODULE=auto GOPROXY=direct go get -u golang.org/x/lint/golint)

.PHONY: info
info:
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"


.PHONY: test
test:
	@echo
	@echo "==> Running unit tests <=="

	GO111MODULE=auto GOPROXY=direct go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: generated-code
generated-code:
	CGO_ENABLED=0 go generate ./...