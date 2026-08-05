.PHONY: init\:local
.PHONY: proto\:gen

init\:local:
	docker network create tr-core-internal-network || true
	docker volume create tr-core-postgres-data || true
	docker volume create tr-core-rabbitmq-data || true
	docker volume create tr-core-redis-data || true

proto\:gen:
	cd proto && buf generate
