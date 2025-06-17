package main

import (
	"log"

	"github.com/t-0-network/provider-sdk-go/service/internal"
)

func main() {
	app, err := internal.NewAppContainer()
	if err != nil {
		log.Fatalf("Failed to create app container: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
