FROM alpine
MAINTAINER Victor Castell <victor@victorcastell.com>

RUN set -x \
	&& buildDeps='bash ca-certificates openssl tzdata' \
	&& apk add --update $buildDeps \
	&& rm -rf /var/cache/apk/* \
	&& mkdir -p /opt/local/dkron

EXPOSE 8080 8946

ENV SHELL /bin/bash
WORKDIR /opt/local/dkron

COPY dist/linux_amd64/. .
ENTRYPOINT ["/opt/local/dkron/dkron"]

CMD ["--help"]
