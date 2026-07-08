run:
	docker compose up --build

stop:
	docker compose down

setup:
	docker compose up -d db api
	sleep 3
	docker compose exec api go mod init api-hotelaria || true
	docker compose exec api go get github.com/joho/godotenv
	docker compose exec api go get github.com/jackc/pgx/v5/pgxpool
	docker compose exec api go get github.com/gin-gonic/gin
	docker compose exec api air init || true
	docker compose down
	@echo "🚀 Setup concluído com sucesso! Agora basta rodar 'make run' para iniciar."

.PHONY: clean - clear database and remove containers
clean:
	docker compose down -v