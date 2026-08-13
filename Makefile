MODULE     := github.com/day0ops/preissuance-ext-authz
BINARY     := preissuance-ext-authz
REPO       ?= australia-southeast1-docker.pkg.dev/field-engineering-apac/kasunt
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DOCKER_TAG := $(shell echo $(VERSION) | sed 's/^v//')
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
PLATFORMS  := linux/amd64,linux/arm64

LDFLAGS := -X $(MODULE)/pkg/version.Version=$(VERSION) \
           -X $(MODULE)/pkg/version.Commit=$(COMMIT) \
           -X $(MODULE)/pkg/version.BuildDate=$(BUILD_DATE)

.PHONY: build
build: bin/$(BINARY)-linux-amd64 bin/$(BINARY)-linux-arm64

bin/$(BINARY)-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $@ .

bin/$(BINARY)-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $@ .

.PHONY: docker-build
docker-build: build
	docker buildx build --platform $(PLATFORMS) -t $(REPO)/$(BINARY):$(DOCKER_TAG) .

.PHONY: docker-push
docker-push: build
	docker buildx build --platform $(PLATFORMS) -t $(REPO)/$(BINARY):$(DOCKER_TAG) --push .

.PHONY: release
release: docker-push

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

.PHONY: clean
clean:
	rm -rf bin/
