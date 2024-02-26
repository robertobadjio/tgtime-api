export GOOS=linux
export GOARCH=amd64

DOCKER_COMPOSE = docker-compose -f docker-compose.yml

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up --build -d

up-recreate:
	$(DOCKER_COMPOSE) up --force-recreate -d

down:
	$(DOCKER_COMPOSE) down

clean:
	$(DOCKER_COMPOSE) down --remove-orphans