FROM golang:1.23.1
LABEL maintainer="Victor Castell <0x@vcastellm.xyz>"

EXPOSE 8080 8946

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go install ./...

CMD ["dkron"]
