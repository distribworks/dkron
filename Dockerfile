FROM golang:1.13
LABEL maintainer="Victor Castell <victor@victorcastell.com>"

EXPOSE 8080 8946
USER root
RUN mkdir -p /app
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN go install ./...
RUN chmod -R 777 /app

# CMD ["dkron","agent","--server"]
