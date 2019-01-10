LINUX_PKGS := $(wildcard dist/*.deb) $(wildcard dist/*.rpm)
.PHONY: fury $(LINUX_PKGS)
fury: $(LINUX_PKGS)
$(LINUX_PKGS):
	fury push $@

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
	cd website; hugo -d ../public
	ghp-import -p public

gen:
	go generate ./dkron/templates
	go generate ./dkron/assets

test:
	@bash --norc -i ./scripts/test
