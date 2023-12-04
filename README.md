# go-ticketos-api
Ticketos is ticket provider for events. Go-ticketos is api for provider.

## Dependencies
- [golang](https://go.dev/) 1.21
- [go-enum](https://github.com/abice/go-enum) v0.5.10
- [wire](https://github.com/google/wire) v0.5.0
- [mockery](https://github.com/vektra/mockery) v2.37.1
- [golangci-lint](https://golangci-lint.run/) v1.55.2
- [docker](https://www.docker.com/) v24.0.7
- [docker compose](https://docs.docker.com/compose/) v2.21.0
- [pre-commit](https://pre-commit.com/) v2.17.0

## Usage
### Tests
- Run tests: `make test`
- Clean test cache `test-clean`

### Setup
- Install git hooks: `pre-commit install --hook-type commit-msg && pre-commit install`
- Check for outdated dependencies and upgrade: `go get -u && go mod tidy`
- Set up env: `make up-utils`
- Clean env: `make clean`
- Stop app with env: `make down`

### Development
- Run linter: `make lint`
- Regenerate wire_gen files: `make wire-regenerate`
- Generate generated files: `make generate`

### Run
- Run app (need to build app firstly) `make run`
- Run app with env via docker compose: `make up`

### Docs 
- openapi
  - path `api/openapi-spec/api.yaml`
  - use `https://editor.swagger.io/` for open configuration in ui