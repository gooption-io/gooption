FROM golang:1.11.5 as builder
# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# build directories
RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

# Go dep!
RUN go get -u github.com/golang/dep/...

# dep ensure
ADD . /go/src/github.com/gooption-io/gooption
WORKDIR /go/src/github.com/gooption-io/gooption
RUN ${GOBIN}/dep ensure -vendor-only
