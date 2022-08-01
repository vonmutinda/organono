include .env

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
	go run cmd/main.go -e .env

hello:
	godo hello

test:
	godo test -- -e .env.test

test-lite:
	godo test-now -- -e .env.test

up:
	docker-compose -f docker-compose.yml up --remove-orphans
	
stop:
	docker-compose stop
