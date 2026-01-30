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

.PHONY: all proto probe controller transfer clean

# build target
all: proto probe controller transfer

proto: $(GO_GEN_FILES)

probe:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=amd64 go build -ldflags=$(BUILD_FLAG) -gcflags="all=-trimpath=$(PWD)" \
				-asmflags="all=-trimpath=$(PWD)" -o $(BUILD_DIR)/$@ cmd/probe/main.go

controller:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=amd64 go build -ldflags=$(BUILD_FLAG) -gcflags="all=-trimpath=$(PWD)" \
				-asmflags="all=-trimpath=$(PWD)" -o $(BUILD_DIR)/$@ cmd/controller/main.go
transfer:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=amd64 go build -ldflags=$(BUILD_FLAG) -gcflags="all=-trimpath=$(PWD)" \
				-asmflags="all=-trimpath=$(PWD)" -o $(BUILD_DIR)/$@ cmd/transfer/main.go

# build protobuf to go
$(GEN_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(GEN_DIR)
	$(PROTOC) $(GO_FLAGS) -I$(PROTO_DIR) $<

clean:
	rm -rf $(GEN_DIR)/*.go
	rm -rf $(BUILD_DIR)/*
