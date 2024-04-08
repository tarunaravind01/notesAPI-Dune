package initializers

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// Load varaibles from .env into ENV
func LoadEnvVariable() {
	err := godotenv.Load()
  	if err != nil {
	fmt.Println(err)
    log.Fatal("Error loading .env file")
  }
}