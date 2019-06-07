FROM golang:alpine AS builder

ADD ./ /tmp/redash-exporter/

RUN apk update && \
    apk add git build-base && \
    rm -rf /var/cache/apk/* && \
    mkdir -p "$GOPATH/src/github.com/buidlsville/" && \
    mv /tmp/redash-exporter "$GOPATH/src/github.com/buidlsville/" && \
    cd "$GOPATH/src/github.com/buidlsville/redash-exporter" && \
    GOOS=linux GOARCH=amd64 go build -o redash-exporter && \
    mv redash-exporter /redash-exporter

FROM alpine:3.7

RUN apk add --update ca-certificates

COPY --from=builder /redash-exporter /redash-exporter

ENTRYPOINT ["/redash-exporter"]
