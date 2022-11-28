default: gen

gen:
	rm -rf codegen
	mkdir -p tmp
	git clone https://github.com/awlsring/ActionRunnerModel.git tmp/ActionRunnerModel
	cd tmp/ActionRunnerModel && make
	openapi-generator generate -i tmp/ActionRunnerModel/build/smithyprojections/ActionRunnerModel/open-api/openapi/ActionRunner.openapi.json -g go -o ./codegen/ActionRunnerGoClient --git-user-id awlsring --git-repo-id action-runner
	rm -rf tmp
	mkdir ./codegen/SurrealDBClient
	cp -R /Users/awlsring/Code/SurrealDBClient ./codegen
	go mod tidy

build:
	docker buildx build --platform linux/amd64,linux/arm64 -t awlsring/action-runner --push .