#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker-compose -f docker-compose.yml

include .env
ENV ?= dev

ifdef ENV
ifneq "$(ENV)" ""
	DOCKER_COMPOSE := $(DOCKER_COMPOSE) -f docker-compose.$(ENV).yml
endif
endif

export GOOS=linux
export GOARCH=amd64

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Start services
	$(DOCKER_COMPOSE) up -d $(SERVICES)

build: ## Build service containers
	$(DOCKER_COMPOSE) build

down: ## Down services
	$(DOCKER_COMPOSE) down

clean: ## Delete all containers
	$(DOCKER_COMPOSE) down --remove-orphans

protoc-build: ## Generate pb.go files
	#protoc --proto_path=api/v1/pb/api --go_out=api/v1/pb/api --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:api/v1/pb/api --go-grpc_opt=paths=source_relative apisvc.proto
	protoc --proto_path=api/v1/pb/api --go_out=api/v1/pb/api --go_opt=paths=source_relative --go-grpc_out=api/v1/pb/api --go-grpc_opt=paths=source_relative apisvc.proto
