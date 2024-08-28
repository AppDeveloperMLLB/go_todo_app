.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build the docker image
	docker build -t sample/gotodo:${DOCKER_TAG} \
		--target deploy ./

build-local: ## Build the docker image without cache
	docker compose build --no-cache

up: ## Start the application
	docker compose up -d

down: ## Stop the application
	docker compose down

logs: ## Show the logs
	docker compose logs -f

ps: ## Show the status of the containers
	docker compose ps

test: 		## Run the tests
	go test -race -shuffle=on ./...

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'