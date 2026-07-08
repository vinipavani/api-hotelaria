run:
	docker compose up --build

stop:
	docker compose down

setup:
	DOCKER_BUILDKIT=0 COMPOSE_BAKE=0 docker compose up -d db api
	sleep 3
	docker compose exec api go mod init api-hotelaria || true
	docker compose exec api go mod download
	docker compose exec api air init || true
	docker compose down
	@echo "🚀 Setup concluído com sucesso! Agora basta rodar 'make run' para iniciar."

.PHONY: clean - clear database and remove containers
clean:
	docker compose down -v