up: docker-up
down: docker-down
stop: docker-stop
start: docker-start
restart: docker-restart
build: docker-build
init: create-env docker-down-clear docker-pull docker-build docker-up

create-env:
	if [ ! -f './.env' ]; then cp ./.env.sample ./.env; else exit 0; fi;

docker-up:
	docker-compose up --detach --remove-orphans

docker-down:
	docker-compose down --remove-orphans

docker-down-clear:
	docker-compose down --volumes --remove-orphans

docker-stop:
	docker-compose stop

docker-start:
	docker-compose start

docker-restart:
	docker-compose restart

docker-pull:
	docker-compose pull

docker-build:
	docker-compose build

shell-pg:
	docker-compose exec postgres bash

.PHONY: init \
	run