# Check for gopath
ifndef GOPATH
    $(error You must provide GOPATH env variable.)
endif

# Initialize gonin
GOBIN=$(GOPATH)/bin

# If not GOBIN in PATH - export
ifneq (,$(findstring $GOBIN,$PATH))
	export PATH=$PATH:$GOBIN
endif

# Recipe vars
GOOS=linux
GOARCH=amd64
BUILD_TARGET=cmd/kraken/kraken.go
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Default build recipe
default: setup fmt clean build

setup:
	@go get -u "github.com/tools/godep"
	dep ensure

# Format tool
.PHONY: fmt
fmt:
	gofmt -w $(GOFMT_FILES)

# CLean build dir
.PHONY: clean
clean:
	go clean
	rm -f kraken

# Build kraken
build: setup
	go build --ldflags '-X main.VERSION=0.1' $(BUILD_TARGET)
