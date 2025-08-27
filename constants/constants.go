package constants

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	customDBPassword = ""
)

func CURRENT_YEAR() string {
	return fmt.Sprint(time.Now().Year())
}

func DEFAULT_PASSWORD() string {
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file, using environment variables from OS")
	}

	return os.Getenv("DEFAULT_PASSWORD")
}

func DEFAULT_DB_PASSWORD() string {
	if customDBPassword == "" {
		err := godotenv.Load()
		if err != nil {
			println("Error loading .env file, using environment variables from OS")
		}

		return os.Getenv("DEFAULT_DB_PASSWORD")
	}

	return customDBPassword
}
