all: deps
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh

deps:
	godep save -t ./...

docs:
	mkdocs gh-deploy --clean

prmd:
	prmd doc --prepend docs/docs/overview.md static/schema.json | sed 's/\<a name\=.*a\>//' > docs/docs/api.md

test:
	@bash --norc -i ./scripts/test.sh
