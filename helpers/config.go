package helpers

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	err := godotenv.Load(filepath.Join("./", ".env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
		debug.PrintStack()
	}

	return os.Getenv(key)
}
