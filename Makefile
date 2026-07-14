run:
	docker compose up --build db api

stop:
	docker compose down

setup:
	DOCKER_BUILDKIT=0 COMPOSE_BAKE=0 docker compose up -d db api
	sleep 3
	docker compose exec api go mod init api-hotelaria || true
	docker compose exec api go mod download
	docker compose exec api air init || true
	@make stop
	@make db-migrate
	@make seed
	@echo "🚀 Setup concluído com sucesso! Agora basta rodar 'make run' para iniciar."

seed:
	docker compose run --rm api go run cmd/seed/main.go

test:
	@echo "🧪 Executando a suíte de testes unitários e de integração no banco isolado..."
	docker compose up -d db_test
	@sleep 3
	docker compose run --rm api go test -v ./...
	docker compose stop db_test

.PHONY: clean - clear database and remove containers
clean:
	docker compose down -v

create-migration:
	@if [ -z "$(name)" ]; then echo "⚠️ Erro: Você precisa passar o nome. Ex: make migration name=my_new_table"; exit 1; fi
	docker run --rm -v $(PWD)/internal/database/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)
	@echo "✅ Arquivos de migração criados com sucesso na pasta internal/database/migrations/!"

db-migrate:
	docker compose run --rm api go run cmd/migrate/main.go up

# 📝 Exemplo de uso: "make db-rollback" ou "make db-rollback steps=4" 
steps ?= 1
db-rollback:
	docker compose run --rm api go run cmd/migrate/main.go down $(steps)

