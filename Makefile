all: test

doc:
	# @$(MAKE) apidoc
	cd website; hugo -d ../public
	ghp-import public

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.2.0.jar convert -i docs/swagger.yaml -f website/content/usage/api -c docs/config.properties

gen:
	go generate ./dkron
	go fmt ./dkron/bindata.go

test:
	@bash --norc -i ./scripts/test.sh

release:
	@$(MAKE) doc
	@goxc -tasks+=publish-github
