FROM golang:1.16
LABEL maintainer="Victor Castell <victor@victorcastell.com>"

EXPOSE 8080 8946

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go install ./...

CMD ["dkron"]
