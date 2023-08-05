GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

BINARY_NAME = lite-smtp-proxy

default: build

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -rf dist/

test:
	$(GOCLEAN) -testcache
	$(GOTEST) -v ./...

release:
	goreleaser release --snapshot --clean

all: default

