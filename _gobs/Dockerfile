FROM golang:1.13.1-alpine as builder
RUN apk add --no-cache build-base git ca-certificates && update-ca-certificates 2>/dev/null || true
COPY . /go/src/github.com/gooption-io/gooption
WORKDIR /go/src/github.com/gooption-io/gooption/gobs
ENV GO111MODULE on
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -ldflags="-X github.com/gooption-io/gooption/v1/initializer.Commit=$(git rev-parse --short HEAD) \
    -X github.com/gooption-io/gooption/v1/initializer.Version=$(git describe --tags) \
    -X github.com/gooption-io/gooption/v1/initializer.BuiltAt=$(date +%FT%T%z) \
    -X github.com/gooption-io/gooption/v1/initializer.BuiltOn=$(hostname)" -o gobs

FROM scratch
LABEL maintainer=seb.pouplin@gmail.com
EXPOSE 8080 8081 50051
COPY --from=builder /go/src/github.com/gooption-io/gooption/gobs/gobs /
# COPY gobs.yml /etc/gobs/gobs.yml
CMD ["./gobs"]