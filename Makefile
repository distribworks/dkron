all: deps
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh

deps:
	godep save -t ./...

doc:
	mkdocs gh-deploy --clean

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.0.0.jar convert -i static/swagger.yaml -f docs/docs/api -c docs/config.properties

test:
	@bash --norc -i ./scripts/test.sh

release:
	@$(MAKE) doc
	@$(MAKE) goxc
