.PONY: all build deps image lint test
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	go build -o ./bin/upupaway

buildimage:
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/upupaway
	docker build --no-cache=true --rm --tag upupaway .
	rm -rf ./bin/

run: ## Run it
	./bin/upupaway -c ./config-dev.yml

deps: ## Install dependencies.
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/Masterminds/glide && glide install

image: ## Build the Docker image.
	docker build .

lint: ## Lint the code
	golint $(CHECK_FILES)

vet: # Vet the code
	go vet $(CHECK_FILES)

test: ## Run tests.
	go test -v $(CHECK_FILES)

pushimage: ## Push image to Docker hub
	docker tag upupaway $(DOCKER_ID_USER)/upupaway:latest
	docker push erans/upupaway
