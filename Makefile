GIT_VERSION := $(shell git describe --abbrev=0 --tags)
GIT_REVISION := $(shell git rev-list -1 HEAD)
DATE := $(shell date +%Y-%m-%dT%H:%M%Sz)

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: *.go test ## Build binary
	go build -ldflags "-s -w -X main.version=${GIT_VERSION} -X main.revision=${GIT_REVISION} -X main.buildDate=${DATE}" -trimpath -o bin/ec2id

test: ## Run test
	go test ./...

clean: ## Remove binary
	rm -f bin/ec2id
