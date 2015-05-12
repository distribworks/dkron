all: deps
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh
