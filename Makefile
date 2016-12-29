all: test

doc:
	@$(MAKE) apidoc
	mkdocs gh-deploy --clean

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.1.0.jar convert -i docs/swagger.yaml -f docs/docs/api -c docs/config.properties

test:
	@bash --norc -i ./scripts/test.sh

release:
	@$(MAKE) doc
	@goxc -tasks+=publish-github
