default: gen

gen:
	openapi-generator generate -i /Users/awlsring/Code/DWSRepos/DWS-ActionRunnerModel/build/smithyprojections/DWS-ActionRunnerModel/open-api/openapi/ActionRunner.openapi.json -g go -o ./gen/ActionRunnerGoClient --git-user-id awlsring --git-repo-id action-runner-model
	mkdir ./gen/SurrealDBClient
	cp -R /Users/awlsring/Code/SurrealDBClient ./gen

build:
	docker buildx build --platform linux/amd64,linux/arm64 -t awlsring/action-runner --push .