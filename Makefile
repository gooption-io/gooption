.PHONY: deps tools install

# PROJECTS are subfolders with project files
PROJECTS := gobs gooption
BUILD_VERSION ?= latest
IMAGE_NAME = "gooption/builder:${BUILD_VERSION}"

all: dist

dist: builder image

tools:
	@go get github.com/golang/dep/cmd/dep

deps: tools
	@rm -rf vendor/
	dep ensure -vendor-only

build:
	@for pkg in $(PROJECTS); do \
		cd $$pkg; \
		make build || exit 1; \
		cd -; echo; echo; \
	done

install:
	@for pkg in $(PROJECTS); do \
		cd $$pkg; \
		make install || exit 1; \
		cd -; echo; echo; \
	done

builder:
	docker build --no-cache -t $(IMAGE_NAME) .

image:
	@for pkg in $(PROJECTS); do \
		cd $$pkg; \
		make image || exit 1; \
		cd -; echo; echo; \
	done
