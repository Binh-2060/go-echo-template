package dotenv

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	log.Println("##### Begin load ENV #####")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
