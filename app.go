package main

import (
	"log"
	"store/delivery"

	"github.com/joho/godotenv"
)

func main() {
	// Load config files
	err := godotenv.Load("./config/config.env")
	if err != nil {
		log.Fatalln("Make sure file \"config.env\" is present on this repository")
	}

	// Run the server
	delivery.Server().Run()
}
