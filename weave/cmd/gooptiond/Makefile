.PHONY: all build test image tf protoc clean dist

BUILD_VERSION ?= manual
BUILD_FLAGS := -ldflags "-X github.com/iov-one/weave.Version=${BUILD_VERSION}"
DOCKER_BUILD_FLAGS := -a -installsuffix cgo
BUILDOUT ?= gooptiond
IMAGE_NAME = "gooption/gooptiond:${BUILD_VERSION}"

all: dist

dist: clean test build image

clean:
	rm -f ${BUILDOUT}

build:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) $(DOCKER_BUILD_FLAGS) -o $(BUILDOUT) .

image:
	docker build --pull -t $(IMAGE_NAME) .

test:
	go test -race ./...

# Test fast
tf:
	go test -short ./...

protoc:
	protoc --gogofaster_out=. -I=. -I=$(GOPATH)/src app/*.proto
