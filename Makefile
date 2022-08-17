include .env

version = v1.0.0
dockerpath = vonmutinda/organono

populate_countries:
	go run scripts/populate_countries.go -e .env

fmt:
	go fmt ./...

migrate:
	goose -dir 'app/db/migrations' postgres ${DATABASE_URL} up

resetdb:
	goose -dir 'app/db/migrations' postgres ${DATABASE_URL} reset

rollback:
	goose -dir 'app/db/migrations' postgres ${DATABASE_URL} down

# e.g make migration name=companies 
migration:
	goose -dir app/db/migrations create $(name) sql

server:
	go run cmd/main.go

hello:
	godo hello

test:
	godo test -- -e .env.test

test-lite:
	godo test-lite -- -e .env.test

up:
	docker-compose -f docker-compose.yml up --remove-orphans

stop:
	docker-compose down --remove-orphans

build:
	docker build -t organono . 
	
tag:
	docker tag organono $(dockerpath):$(version)

push:
	docker push vonmutinda/organono:$(version)

run:
	docker run --name organono-$(version) -p ${PORT}:${PORT} organono
