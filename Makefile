include .env_local
export

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run:  ### run server
	go mod tidy & go mod download &&\
	go run ./cmd/gophermart/main.go
.PHONY: run


compose-up-db:  ### up DB
	podman-compose up postgres -d
.PHONY: compose-up-db


up:  ### Run in container
	podman-compose up --build
.PHONY: up
