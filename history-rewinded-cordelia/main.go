package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {



	log.Println("Starting Cordelia...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load enviroment variables, reason: %s", err)
	}


}
