tests: docker
	docker-compose -f docker-compose.yml up

backend:
	docker-compose -f docker-compose-backend.yml up

docker:
	docker build -t f3c .
