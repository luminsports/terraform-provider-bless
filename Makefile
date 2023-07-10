export GO111MODULE=on
VERSION=$(shell cat VERSION)
export BASE_BINARY_NAME=terraform-provider-bless_v$(VERSION)
SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
export DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
LDFLAGS=-ldflags "-w -s -X github.com/luminsports/terraform-provider-bless/pkg/version.GitSha=${SHA} -X github.com/luminsports/terraform-provider-bless/pkg/version.Version=${VERSION} -X github.com/luminsports/terraform-provider-bless/pkg/version.Dirty=${DIRTY}"


setup: ## setup development dependencies
	./.godownloader-packr.sh -d v1.24.1
	curl -sfL https://raw.githubusercontent.com/chanzuckerberg/bff/main/download.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
.PHONY: setup

lint-all: ## run all the linters
	./bin/golangci-lint run
.PHONY: lint-all

build: packr
	@CGO_ENABLED=0 go build ${LDFLAGS} -o $(BASE_BINARY_NAME) .
.PHONY:  build

test: deps packr
	@TF_ACC=yes go test -cover -v ./...
.PHONY: test

test-ci: packr
	@TF_ACC=yes go test -cover -v ./...
.PHONY: test-ci

deps:
	go mod tidy
.PHONY: deps

packr: clean
	./bin/packr
.PHONY: packr

clean: ## clean the repo
	rm terraform-provider-bless 2>/dev/null || true
	go clean
	rm -rf dist 2>/dev/null || true
	./bin/packr clean
	rm coverage.out 2>/dev/null || true
.PHONY: clean

check-release-prereqs:
ifndef KEYBASE_KEY_ID
	$(error KEYBASE_KEY_ID is undefined)
endif
.PHONY: check-release-prereqs

release: check-release-prereqs ## run a release
	./bin/bff bump
	git push
	goreleaser release --debug --rm-dist
.PHONY: release

install-tf: build ## installs plugin where terraform can find it
	mkdir -p $(HOME)/.terraform.d/plugins
	cp ./$(BASE_BINARY_NAME) $(HOME)/.terraform.d/plugins/$(BASE_BINARY_NAME)
.PHONY: install-tf
