package main

import (
	"kelp-go-assessment/internal/api/routes"
	"log"
)

func main() {
	r := routes.SetupRouter()

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
