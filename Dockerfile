FROM golang:1.7.3
MAINTAINER Victor Castell <victor@victorcastell.com>

EXPOSE 8080 8946

RUN wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz -O /tmp/glide.tar.gz && \
    tar zxvf /tmp/glide.tar.gz -C /tmp && \
    mv /tmp/linux-amd64/glide /usr/local/bin/ && \
    rm -rf /tmp/glide.tar.gz /tmp/linux-amd64

RUN mkdir -p /gopath/src/github.com/victorcoder/dkron
WORKDIR /gopath/src/github.com/victorcoder/dkron

ENV GOPATH /gopath
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

COPY glide.yaml ./glide.yaml
COPY glide.lock ./glide.lock
RUN glide install

COPY . ./
RUN go build *.go
CMD ["/gopath/src/github.com/victorcoder/dkron/main"]
