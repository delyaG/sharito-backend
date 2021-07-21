PROJECT_NAME := sharito
PROJECT := gitlab.com/sharito/backend.git

VERSION := $(shell cat version)
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := "-s -w -X $(PROJECT)/internal/version.Version=$(VERSION) -X $(PROJECT)/internal/version.Commit=$(COMMIT)"

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)

test:
	@go test -v -cover -gcflags=-l --race ./...

GOLANGCI_LINT_VERSION := v1.31.0
lint:
	@golangci-lint run -v

dep:
	@go mod download

docker-run:
	docker run --name=$(CONTAINER_NAME) -p $(PORT):8080 \
		--env-file $(ENV_FILE) \
		--network=prod \
		-e SHARITO_POSTGRES_PASSWORD="$(SHARITO_POSTGRES_PASSWORD)" \
		-v $(SHARITO_HTTP_JWT_PRIVATE_KEY):/jwt.key \
		--restart always -d $(CI_REGISTRY_IMAGE):latest

docker-stop:
	docker stop $(CONTAINER_NAME) || true

docker-rm:
	docker rm -f $(CONTAINER_NAME) || true

docker-rerun: docker-stop docker-rm docker-run
