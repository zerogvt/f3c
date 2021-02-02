tests: docker
	docker-compose -f docker-compose-tests.yml up

backend:
	docker-compose -f 

docker:
	docker build -t f3c .
