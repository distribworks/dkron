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

.PHONY: doc apidoc gen test
doc:
	#scripts/run doc --dir website/content/cli
	cd website; hugo -d ../public
	ghp-import -p public

gen:
	rm -rf static/.node_modules
	go generate ./dkron/templates
	go generate ./dkron/assets

test:
	@bash --norc -i ./scripts/test

updatetestcert:
	wget https://badssl.com/certs/badssl.com-client.p12 -q -O badssl.com-client.p12
	openssl pkcs12 -in badssl.com-client.p12 -nocerts -nodes -passin pass:badssl.com -out builtin/bins/dkron-executor-http/testdata/badssl.com-client-key-decrypted.pem
	openssl pkcs12 -in badssl.com-client.p12 -nokeys -passin pass:badssl.com -out builtin/bins/dkron-executor-http/testdata/badssl.com-client.pem
	rm badssl.com-client.p12
