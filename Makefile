default: gen

gen:
	rm -rf codegen
	mkdir -p tmp
	mkdir -p codegen
	git clone https://github.com/awlsring/ActionRunnerModel.git tmp/ActionRunnerModel
	git clone https://github.com/awlsring/SurrealDBClient.git tmp/SurrealDBClient
	mkdir ./codegen/SurrealDBClient
	cp -R tmp/SurrealDBClient ./codegen
	cd tmp/ActionRunnerModel && gradle build
	openapi-generator generate -i tmp/ActionRunnerModel/build/smithyprojections/ActionRunnerModel/open-api/openapi/ActionRunner.openapi.json -g go -o ./codegen/ActionRunnerGoClient --git-user-id awlsring --git-repo-id action-runner
	rm -rf tmp
	go mod tidy

build:
	docker buildx build --platform linux/amd64,linux/arm64 -t awlsring/action-runner --push .