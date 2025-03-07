FROM golang:1.23.1
LABEL maintainer="Victor Castell <0x@vcastellm.xyz>"

EXPOSE 8080 8946

RUN mkdir -p /app
WORKDIR /app

ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on

# Leverage build cache by copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN go mod verify

RUN go mod download

# Copy the rest of the source code
COPY . .

RUN go install ./...

CMD ["dkron"]
