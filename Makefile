LINUX_PKGS := $(wildcard dist/*.deb) $(wildcard dist/*.rpm)
.PHONY: fury $(LINUX_PKGS)
fury: $(LINUX_PKGS)
$(LINUX_PKGS):
	fury push --as distribworks $@

.PHONY: goreleaser
goreleaser:
	docker run --rm --privileged \
	-v ${PWD}:/dkron \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-w /dkron \
	-e GITHUB_TOKEN \
	-e DOCKER_USERNAME \
	-e DOCKER_PASSWORD \
	-e DOCKER_REGISTRY \
	--entrypoint "" \
	goreleaser/goreleaser scripts/release

.PHONY: release
release: clean goreleaser

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

.PHONY: doc apidoc test ui updatetestcert
doc:
	#scripts/run doc --dir website/content/cli
	cd website; hugo -d ../public
	ghp-import -p public

test:
	@bash --norc -i ./scripts/test

localtest:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

updatetestcert:
	wget https://badssl.com/certs/badssl.com-client.p12 -q -O badssl.com-client.p12
	openssl pkcs12 -in badssl.com-client.p12 -nocerts -nodes -passin pass:badssl.com -out builtin/bins/dkron-executor-http/testdata/badssl.com-client-key-decrypted.pem
	openssl pkcs12 -in badssl.com-client.p12 -nokeys -passin pass:badssl.com -out builtin/bins/dkron-executor-http/testdata/badssl.com-client.pem
	rm badssl.com-client.p12

ui/node_modules: ui/package.json
	cd ui; npm install
	# touch the directory so Make understands it is up to date
	touch ui/node_modules

dkron/ui-dist: ui/node_modules ui/public/* ui/src/* ui/src/*/*
	cd ui; npm run-script build

plugin/types/%.pb.go: proto/%.proto
	protoc -I proto/ --go_out=plugin/types --go_opt=paths=source_relative --go-grpc_out=plugin/types --go-grpc_opt=paths=source_relative $<

ui: dkron/ui-dist

main: dkron/ui-dist plugin/types/dkron.pb.go plugin/types/executor.pb.go *.go */*.go */*/*.go */*/*/*.go
	GOBIN=`pwd` go install ./builtin/...
	go mod tidy
	go build main.go
