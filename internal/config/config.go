package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port            string
	DatabaseURL     string
	TestDatabaseURL string
}

var Env *Config

func LoadConfig() {
	port := getEnvVariable("PORT", "8080")
	databaseURL := getEnvVariable("DATABASE_URL", "")
	TestDatabaseURL := getEnvVariable("TEST_DATABASE_URL", "")

	Env = &Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		TestDatabaseURL: TestDatabaseURL,
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
