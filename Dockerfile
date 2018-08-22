FROM golang:1.11-rc
MAINTAINER Victor Castell <victor@victorcastell.com>

EXPOSE 8080 8946

RUN mkdir -p /app
WORKDIR /app

ENV GO111MODULE=on
COPY . .
RUN go install ./...
CMD ["dkron"]
