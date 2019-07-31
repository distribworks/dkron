FROM golang:1.12-alpine
LABEL maintainer="Victor Castell <victor@victorcastell.com>"

EXPOSE 8080 8946

RUN apk add --no-cache git mercurial bash

RUN mkdir -p /app
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go install ./...

CMD ["dkron"]
