LINUX_PKGS := $(wildcard dist/*.deb) $(wildcard dist/*.rpm)
.PHONY: fury $(LINUX_PKGS)
fury: $(LINUX_PKGS)
$(LINUX_PKGS):
	fury push $@

.PHONY: goreleaser
goreleaser:
	goreleaser --rm-dist

.PHONY: release
release: clean goreleaser fury

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
	go generate ./dkron
	go fmt ./dkron/bindata.go

test:
	@bash --norc -i ./scripts/test
