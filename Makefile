include .env

fmt:
	go fmt ./...

migrate:
	goose -dir 'app/db/migrations' postgres ${DATABASE_URL} up

resetdb:
	goose -dir 'app/db/migrations' postgres ${DATABASE_URL} reset

rollback:
	goose -dir 'app/db/db/migrations' postgres ${DATABASE_URL} down

# e.g make migration name=companies 
migration:
	goose -dir app/db/migrations create $(name) sql

server:
	go run cmd/main.go -e .env

test:
	export DATABASE_URL="postgres://admin:password@localhost:5433/organono_backend?sslmode=disable" && \
	go test ./...

up:
	docker-compose -f docker-compose.yml up --remove-orphans
	
stop:
	docker-compose stop
