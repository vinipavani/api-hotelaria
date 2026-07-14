package config

import (
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port        string
	DatabaseURL string
}

var Env *Config

func LoadEnv() {
	port := getEnvVariable("PORT", "8080")

	var databaseURL string
	if isTestEnvironment() {
		databaseURL = getEnvVariable("TEST_DATABASE_URL", "")
	} else {
		databaseURL = getEnvVariable("DATABASE_URL", "")
	}

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

func isTestEnvironment() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") || strings.Contains(arg, "test") {
			return true
		}
	}
	return false
}
