all: test

doc:
	mkdocs gh-deploy --clean

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.0.0.jar convert -i static/swagger.yaml -f docs/docs/api -c docs/config.properties

test:
	@bash --norc -i ./scripts/test.sh

release:
	@$(MAKE) apidoc
	@$(MAKE) doc
	@goxc -tasks+=publish-github
