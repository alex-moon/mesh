run:
	docker-compose -f docker-compose.local.yml up

build:
	docker-compose -f docker-compose.local.yml build --progress plain

deploy:
	./deploy.sh
