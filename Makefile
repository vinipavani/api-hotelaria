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

create-migration:
	@if [ -z "$(name)" ]; then echo "⚠️ Erro: Você precisa passar o nome. Ex: make migration name=my_new_table"; exit 1; fi
	docker run --rm -v $(PWD)/internal/database/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)
	@echo "✅ Arquivos de migração criados com sucesso na pasta internal/database/migrations/!"