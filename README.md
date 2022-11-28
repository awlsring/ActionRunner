# ActionRunner

## What

ActionRunner is an API layer that sit on top of Ansible to create executions of specified playbooks on given hosts. This project is still in early stages and will be changed over time.

ActionRunner is defined via a Smithy model, which can be used to generate Smithy or OpenAPI clients to call the service. The model package can be found on the [ActionRunnerModel](https://github.com/awlsring/ActionRunnerModel) repo.

ActionRunner uses a given playbook directory to register as targets for executions. This can be any Ansible playbook directory. The service will watch this directory for updates and will register new playbooks as targets. An example of playbooks can be found at my [ActionRunnerPlaybooks](https://github.com/awlsring/ActionRunnerPlaybooks) repo.

## Setup

An example deployment using docker compose can be found under examples/docker-compose

A SurrealDB instance is required also to save executions.

A config.yaml file will need to be passed, an example of this looks like

```yaml
logLevel: debug
api:
  port: "7032"
runner:
  ansibleUser: action-runner
  connectionType: ssh
  playbookSource: "../ActionRunnerPlaybooks"
  privateKeyFile: "../ActionRunnerPlaybooks/keys/id_rsa"
db:
  surreal:
    user: root
    password: root
    address: "ws://localhost:8000/rpc"
    namespace: action-runner
    database: action-runner
```