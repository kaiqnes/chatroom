# SHELL=/bin/bash
# .ONESHELL:
# .DEFAULT_GOAL := all

SERVER_SOURCES:=./cmd
SERVER_TARGET:=./chatroom

.ONESHELL:
.PHONY: install-tools
install-tools: install-asdf-tools install-go-tools asdf-reshim

.ONESHELL:
.PHONY: install-asdf-tools
install-asdf-tools:
	@cat .tool-versions | awk '{print $$1}' | xargs -L 1 asdf plugin add; \
	asdf install

.ONESHELL:
.PHONY: install-go-tools
install-go-tools:
	go install github.com/swaggo/swag/cmd/swag@latest; \
	go install github.com/vektra/mockery/v2@v2.23.4; \

.ONESHELL:
.PHONY: asdf-reshim
asdf-reshim:
	asdf reshim golang && export GOROOT="$(asdf where golang)/go/" && export GOPATH="$(asdf where golang)/packages/";

.PHONY: download
download: install-tools
	@go mod download

.PHONY: format
format:
	goimports -w -l .
	golangci-lint run --fix

.PHONY: check-formatting
check-formatting:
	test -z "$(shell goimports -l .)"

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: run-db-local
run-db-local: stop-local
	sleep 5 && docker-compose up -d

.PHONY: stop-local
stop-local:
	docker-compose down

.PHONY: server-build
server-build:
	go build \
	-ldflags \
    $(LD_FLAGS) \
	-o "$(SERVER_TARGET)" $(SERVER_SOURCES)

.PHONY: run-local
run-local: run-db-local server-build
	ENV=dev $(SERVER_TARGET)

.PHONY: migrate-up
migrate-up:
	./scripts/migrate.sh $(ENV) up

.PHONY: migrate-create
migrate-create:
	./scripts/migrate-create.sh $(name)

.PHONY: migrate-down
migrate-down:
	./scripts/migrate.sh $(ENV) down
