package config

import(
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

var Env *Config

func LoadConfig() {
	port := getEnvVariable("PORT", "8080")
	databaseURL := getEnvVariable("DATABASE_URL", "")

	Env = &Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}
}

func getEnvVariable(key string, defaultValue string) string {
	variable := os.Getenv(key)
	if variable == "" {
		if defaultValue == "" {
			log.Fatalf("A variável de ambiente %s é obrigatória e não foi definida.", key)
		}
		return defaultValue
	}

	return variable
}