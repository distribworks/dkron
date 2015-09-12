all: deps
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh

deps:
	go get -d -v ./...

prmd:
	prmd doc --prepend docs/docs/overview.md static/schema.json > docs/docs/api.md

