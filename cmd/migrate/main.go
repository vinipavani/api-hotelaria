package main

import (
	"api-hotelaria/internal/config"
	"api-hotelaria/internal/database"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	action := os.Args[1]
	config.LoadConfig()
	log.Println("🔌 Conectando ao banco de dados para gerenciar migrações administrativas...")
	database.ConnectDB()
	defer database.CloseDB()

	migration := database.MigrationInstance()
	defer migration.Close()

	switch action {
	case "up":
		log.Println("🚀 Varrendo pacotes embutidos e aplicando novas migrações (UP)...")
		if err := migration.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("✨ O banco de dados já está na última versão. Nenhuma ação necessária.")
				return
			}
			log.Fatalf("❌ Falha crítica ao aplicar migrações: %v", err)
		}
		log.Println("✅ Estrutura do banco de dados erguida e atualizada com sucesso!")

	case "down":
		stepsToRollback := -1

		if len(os.Args) >= 3 {
			customSteps, err := strconv.Atoi(os.Args[2])
			if err != nil || customSteps <= 0 {
				log.Fatalf("❌ Erro: O número de passos para o rollback deve ser um valor inteiro positivo. Recebeu: '%s'", os.Args[2])
			}

			stepsToRollback = customSteps * -1
		}

		log.Printf("🔄 Executando Rollback estrito de exatamente %d versão(ões) aplicada(s) (DOWN)... \n", stepsToRollback*-1)

		if err := migration.Steps(stepsToRollback); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("✨ Nenhuma migração ativa encontrada para sofrer rollback.")
				return
			}
			log.Fatalf("❌ Falha crítica ao desfazer a(s) migração(ões): %v", err)
		}
		log.Println("✅ Rollback executado e tabelas removidas com sucesso!")

	default:
		log.Fatalf("❌ Comando inválido: '%s'. Use apenas 'up' ou 'down'.", action)
	}
}
