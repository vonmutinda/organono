include .env

version = latest
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
	docker run --name organono-$(version) -p ${PORT}:${PORT} --env-file=.env organono

rds-setup:
	aws cloudformation create-stack --stack-name organono-rds --template-body file://infra/rds.yml --parameters file://infra/rds_params.json --capabilities "CAPABILITY_IAM" "CAPABILITY_NAMED_IAM" --region=us-east-1 --profile=default

delete-rds:
	aws cloudformation delete-stack --stack-name organono-rds --region=us-east-1 --profile=default

create-cluster:
	eksctl create cluster --name organono-cluster --region=us-east-1 --nodes=2 --profile=default

delete-cluster:
	eksctl delete cluster --name organono-cluster --region=us-east-1 --profile=default

deploy:
	kubectl apply -f ./infra/pods/organono-api-deployment.yml

cluster-status:
	kubectl get deploy,rs,svc,pods

printenv:
	kubectl exec organono-api -- printenv

logs:
	kubectl logs pod/organono-api

port-forward:
	kubectl port-forward pod/organono-api ${PORT}:80
