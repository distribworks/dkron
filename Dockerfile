FROM golang:1.6-alpine
MAINTAINER Victor Castell <victor@victorcastell.com>

EXPOSE 8080 8946

RUN set -x \
	&& buildDeps='bash git' \
	&& apk add --update $buildDeps \
	&& rm -rf /var/cache/apk/*

RUN mkdir -p /gopath/src/github.com/victorcoder/dkron
WORKDIR /gopath/src/github.com/victorcoder/dkron

ENV GO15VENDOREXPERIMENT=1
ENV GOPATH=/gopath

CMD go run *.go
