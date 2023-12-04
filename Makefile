BINDIR := $(CURDIR)/bin
API_BIN_NAME := api

# go option
PKG        := ./...
TAGS       :=
TESTFLAGS  := -p 1
LDFLAGS    := -w -s
GOFLAGS    :=


build:
	CGO_ENABLED=0 go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(API_BIN_NAME)' ./cmd/$(API_BIN_NAME)

run:
	$(BINDIR)/$(API_BIN_NAME)

test:
	ENV=tests go test $(GOFLAGS) $(PKG) $(TESTFLAGS)

test-clean:
	go clean -testcache

lint:
	golangci-lint run

up:
	docker compose -f docker-compose.yaml up --always-recreate-deps --force-recreate -d --build

down:
	docker compose -f docker-compose.yaml down

up-utils:
	docker compose -f docker-compose.yaml up --always-recreate-deps --force-recreate -d --build database

clean: down
	docker volume rm go-ticketos_pg_data

generate:
	go generate ./...

wire-regenerate:
	rm ./cmd/*/wire_gen.go || true
	go run github.com/google/wire/cmd/wire@v0.5.0 ./cmd/...
