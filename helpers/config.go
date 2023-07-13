package helpers

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	path, err := filepath.Abs(".env")

	if err != nil {
		debug.PrintStack()
		log.Fatalf("Could not find .env file")
	}

	err = godotenv.Load(path)

	if err != nil {
		debug.PrintStack()
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
