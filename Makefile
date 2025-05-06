LINUX_PKGS := $(wildcard dist/*.deb) $(wildcard dist/*.rpm)
.PHONY: fury $(LINUX_PKGS)
fury: $(LINUX_PKGS)
$(LINUX_PKGS):
	fury push --as distribworks $@

PACKAGE_NAME          := github.com/distribworks/dkron
GOLANG_CROSS_VERSION  ?= v1.22

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-v ${PWD}:/dkron \
		-w /dkron \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e DOCKER_REGISTRY \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate,publish --timeout=1h --parallelism=1

.PHONY: release
release:
	@docker run \
		--rm \
		--privileged \
		-v ${PWD}:/dkron \
		-w /dkron \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e DOCKER_REGISTRY \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate --timeout=1h

.PHONY: clean
clean:
	rm -f main
	rm -f *_SHA256SUMS
	rm -f dkron-*
	rm -rf build/*
	rm -rf builder/skel/*
	rm -f *.deb
	rm -f *.rpm
	rm -f *.tar.gz
	rm -rf tmp
	rm -rf ui-dist
	rm -rf ui/build
	rm -rf ui/node_modules
	GOBIN=`pwd` go clean -i ./builtin/...
	GOBIN=`pwd` go clean

.PHONY: docs apidoc test ui updatetestcert
docs:
	# scripts/run doc --dir website/docs/cli
	
	# Build with docker while bun reach compatibility with docusaurs
	cd website; yarn build --out-dir ../public
	ghp-import -p public

test:
	@bash --norc -i ./scripts/test

localtest:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

updatetestcert:
	wget https://badssl.com/certs/badssl.com-client.p12 -q -O badssl.com-client.p12
	openssl pkcs12 -in badssl.com-client.p12 -nocerts -nodes -passin pass:badssl.com -out plugin/http/testdata/badssl.com-client-key-decrypted.pem
	openssl pkcs12 -in badssl.com-client.p12 -nokeys -passin pass:badssl.com -out plugin/http/testdata/badssl.com-client.pem
	rm badssl.com-client.p12

ui/node_modules: ui/package.json
	cd ui; bun install
	# touch the directory so Make understands it is up to date
	touch ui/node_modules

dkron/ui-dist: ui/node_modules ui/public/* ui/src/* ui/src/*/*
	rm -rf dkron/ui-dist
	cd ui; yarn build --out-dir ../dkron/ui-dist

# deprecated
# proto: types/dkron.pb.go types/executor.pb.go types/pro.pb.go
.PHONY: proto
proto:
	@echo "Generating protobuf files using Buf"
	buf generate

lint: proto-lint

proto-lint:
	@echo "Linting protobuf files using Buf"
	buf lint

# deprecated
types/%.pb.go: proto/%.proto
	protoc -I proto/ --go_out=types --go_opt=paths=source_relative --go-grpc_out=types --go-grpc_opt=paths=source_relative $<

client:
	oapi-codegen -package client website/static/openapi/openapi.yaml > client/client.go

ui: dkron/ui-dist

main: dkron/ui-dist types/dkron.pb.go types/executor.pb.go *.go */*.go */*/*.go */*/*/*.go
	GOBIN=`pwd` go install ./builtin/...
	go mod tidy
	go build main.go
