package main

import (
	"log"
	"stride-wars-app/internal/repository"
)

func main() {
	log.Println("Running server...")

	client, err := repository.InitEnt()
	if err != nil {
		log.Fatalf("Error during Ent initialization: %v", err)
	}
	defer client.Close()
}
