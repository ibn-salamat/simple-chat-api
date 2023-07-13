package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	entries, err := os.ReadDir("../")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}

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
