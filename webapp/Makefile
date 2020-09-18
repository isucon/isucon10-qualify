SHELL=/bin/bash

include ../bench/Makefile
.DEFAULT_GOAL := help

api-server/%: ## docker-compose up with mysql and api-server
	docker-compose down -v
	docker-compose -f docker-compose/$(shell basename $@).yaml up --build mysql api-server

isuumo/%: ## docker-compose up with mysql and api-server
	docker-compose down -v
	docker-compose -f docker-compose/$(shell basename $@).yaml up --build mysql api-server nginx frontend

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+/?[a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
