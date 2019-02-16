.PHONY: all build test tf clean image dist

BUILD_VERSION ?= manual
BUILD_FLAGS := -a -tags netgo -tags nocgo -ldflags '-extldflags "-static" -X "github.com/gooption-io/gooption/gobs.Version=$(GIT_VERSION)"'
BUILDOUT ?= gobs
IMAGE_NAME = "gooption/gobs:${BUILD_VERSION}"

all: dist

dist: clean test build image

clean:
	rm -f ${BUILDOUT}

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILDOUT) .

test:
	go test -race ./...

# Test fast
tf:
	go test -short ./...

image:
	docker build --no-cache -t $(IMAGE_NAME) .
