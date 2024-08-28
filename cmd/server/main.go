package main

import (
	"log"
	"net/http"
	"os"

	"github.com/satioO/fhir/v2/router"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
		log.Printf("defaulting to port %s", port)
	}

	handlers := router.RegisterRoutes()
	log.Fatal(http.ListenAndServe(":"+port, handlers))
}
