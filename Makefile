BIN_NAME := dkron
DEP_VERSION=0.4.1
VERSION := 0.10.3
PKGNAME := dkron
LICENSE := LGPL 3.0
VENDOR=
URL := https://dkron.io
RELEASE := 0
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := amd64
PLUGINS=$(wildcard builtin/bins/dkron-*)
BINS=$(wildcard build/*/dkron-*)
DESC := Distributed, fault tolerant job scheduling system
MAINTAINER := Victor Castell <victor@distrib.works>
DOCKER_WDIR := /tmp/fpm
DOCKER_FPM := devopsfaith/fpm

PLATFORMS := darwin_amd64 linux_amd64 linux_arm windows_amd64
TEMP = $(subst _, ,$@)
GOOS = $(word 1, $(TEMP))
GOARCH = $(word 2, $(TEMP))
TARS := $(foreach t,$(PLATFORMS),builder/skel/$(t))

LDFLAGS=-ldflags="-X github.com/victorcoder/dkron/dkron.Version=${VERSION}"

FPM_OPTS=-s dir -v $(VERSION) -n $(PKGNAME) \
  --license "$(LICENSE)" \
  --vendor "$(VENDOR)" \
  --maintainer "$(MAINTAINER)" \
  --architecture $(ARCH) \
  --url "$(URL)" \
  --description  "$(DESC)" \
	--config-files etc/ \
  --verbose

DEB_OPTS=
RPM_OPTS= #--rpm-sign

default: build

all: clean release

release: deb rpm tgz
	
.PHONY: build
build:
	$(foreach p,$(PLUGINS),$(shell go build -o $(shell basename $(p)) ./$(p)))
	go build ${LDFLAGS} -o main .

build_all: $(PLATFORMS)
$(PLATFORMS):
	$(eval ext := $(if $(filter $(GOOS),windows),.exe))
	$(foreach p,$(PLUGINS),$(shell GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/$(GOOS)_$(GOARCH)/$(shell basename $(p))$(ext) ./$(p)))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build ${LDFLAGS} -o build/$(GOOS)_$(GOARCH)/${BIN_NAME}$(ext) .

builder/skel/%/etc/dkron/dkron.yml: builder/files/dkron.yml
	mkdir -p "$(dir $@)"
	cp $< "$@"

builder/skel/%/lib/systemd/system/dkron.service: builder/files/dkron.service
	mkdir -p "$(dir $@)"
	cp builder/files/dkron.service "$@"

builder/skel/%/usr/bin: build_all
	mkdir -p $@
	cp build/linux_amd64/* $@

builder/skel/%: build_all
	mkdir -p $@
	cp README.md LICENSE builder/files/dkron.yml $(wildcard build/$*)/* $@

.PHONY: tgz
tgz: $(addprefix builder/skel/,${PLATFORMS})
	$(foreach p,$(PLATFORMS),$(shell tar zcvf dkron_${VERSION}_${p}.tar.gz -C builder/skel/${p} .))

.PHONY: deb
deb: builder/skel/deb/usr/bin
deb: builder/skel/deb/etc/dkron/dkron.yml
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE} \
		--deb-systemd builder/files/dkron.service \
		--chdir builder/skel/$@ \
		${FPM_OPTS}
	docker build --build-arg debfile=dkron_${VERSION}-${RELEASE}_amd64.deb -f tests/deb/Dockerfile -t test_dkron_${VERSION}_deb .

.PHONY: rpm
rpm: builder/skel/rpm/usr/bin
rpm: builder/skel/rpm/usr/lib/systemd/system/dkron.service
rpm: builder/skel/rpm/etc/dkron/dkron.yml
	docker run --rm -it -v "${PWD}/rpmmacros:/root/.rpmmacros" \
		-v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:rpm -t rpm ${RPM_OPTS} \
		--iteration ${RELEASE} \
		-C builder/skel/$@ \
		${FPM_OPTS}
	docker build --build-arg rpmfile=dkron-${VERSION}-${RELEASE}.x86_64.rpm -f tests/rpm/Dockerfile -t test_dkron_${VERSION}_rpm .

PKGS := $(wildcard *.tar.gz) $(wildcard *.deb) $(wildcard *.rpm)
.PHONY: ghrelease $(PKGS) github
ghrelease:
	github-release release \
		--user victorcoder \
		--repo dkron \
		--tag v${VERSION} \
		--name "${VERSION}" \
		--description "See: https://github.com/victorcoder/dkron/blob/master/CHANGELOG.md" \

$(PKGS): ghrelease
		github-release upload \
			--user victorcoder \
			--repo dkron \
			--name $@ \
			--tag v${VERSION} \
			--file $@

github: $(PKGS)

.PHONY: clean
clean:
	rm -f main
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

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.2.0.jar convert -i website/content/swagger.yaml -f website/content/usage/api -c docs/config.properties

gen:
	go generate ./dkron
	go fmt ./dkron/bindata.go

test:
	@bash --norc -i ./scripts/test.sh
