.PHONY: init\:local

init\:local:
	-docker network create tr-core-internal-network
	-docker volume create tr-core-postgres-data
	-docker volume create tr-core-rabbitmq-data
	-docker volume create tr-core-redis-data
