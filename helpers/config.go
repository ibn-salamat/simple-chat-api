package helpers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	err := godotenv.Load(filepath.Join("./", ".env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
