package main

import (
	"log"
	"os"

	"github.com/nurcholisnanda/hotel-management-system/routes"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	r := routes.SetupRouter()
	r.Run(":" + port)
}
