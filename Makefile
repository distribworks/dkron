all: test

doc:
	@$(MAKE) apidoc
	hugo

apidoc:
	java -jar ~/bin/swagger2markup-cli-1.2.0.jar convert -i docs/swagger.yaml -f docs/docs/api -c docs/config.properties

gen:
	go generate ./dkron
	go fmt ./dkron/bindata.go

test:
	@bash --norc -i ./scripts/test.sh

release:
	@$(MAKE) doc
	@goxc -tasks+=publish-github
