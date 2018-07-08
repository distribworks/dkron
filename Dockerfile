FROM golang:1.10
MAINTAINER Victor Castell <victor@victorcastell.com>

EXPOSE 8080 8946

RUN wget https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/victorcoder/dkron
WORKDIR /go/src/github.com/victorcoder/dkron

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY . .
RUN go install ./...
CMD ["dkron"]
