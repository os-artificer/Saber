# base
PROJECT  := saber
PROTOC    = protoc
PROTO_DIR = pkg/proto/idl
GEN_DIR   = pkg/proto
BUILD_DIR = build
GO_OS    ?= linux
MODULE    = os-artificer/saber/pkg

# version
BUILDTIME  = $(shell date +%Y-%m-%dT%T%z)
GITTAG     = $(shell git describe --tags --always)
GITHASH    = $(shell git rev-parse --short HEAD)
VERSION   ?= $(GITTAG)-$(shell date +%y.%m.%d)

BUILD_FLAG = "-X '$(MODULE)/version.buildTime=$(BUILDTIME)' \
			 -X '$(MODULE)/version.gitTag=$(GITTAG)' \
			 -X '$(MODULE)/version.gitHash=$(GITHASH)' \
             -X '$(MODULE)/version.version=$(VERSION)' "

# flags
GO_FLAGS = --go_out=$(GEN_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative

# search .proto file
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

# generate go code files
GO_GEN_FILES=$(PROTO_FILES:$(PROTO_DIR)/%.proto=$(GEN_DIR)/%.pb.go)

# binaries: one target per name, built from cmd/$(name)/main.go
BINARIES := admin agent controller databus

.PHONY: all proto docker-build clean $(BINARIES)

all: proto $(BINARIES)

proto: $(GO_GEN_FILES)

$(BINARIES):
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=$(GO_OS) GOARCH=amd64 go build -ldflags=$(BUILD_FLAG) \
		-gcflags="all=-trimpath=$(PWD)" -asmflags="all=-trimpath=$(PWD)" \
		-o $(BUILD_DIR)/saber-$@ cmd/$@/main.go

# docker image (build from repo root: make docker-build)
REGISTRY ?= saber
DOCKER_IMAGE = $(REGISTRY)/saber:$(VERSION)
docker-build:
	docker build -f deploy/docker/Dockerfile -t $(DOCKER_IMAGE) .
	@echo "Built $(DOCKER_IMAGE)"

# build protobuf to go
$(GEN_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(GEN_DIR)
	$(PROTOC) $(GO_FLAGS) -I$(PROTO_DIR) $<

clean:
	rm -rf $(GEN_DIR)/*.go
	rm -rf $(BUILD_DIR)/*
