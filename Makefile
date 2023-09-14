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

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: stop-local
stop-local:
	docker-compose down

.PHONY: server-build
server-build:
	go build ./cmd/main.go

.PHONY: run-local
run-local:
	docker-compose up

.PHONY: migrate-create
migrate-create:
	./scripts/migrate-create.sh $(name)

