package main

import (
	"log"
	"os"

	"github.com/openlucky/openlucky/internal/app"
)

func main() {
	configPath := os.Getenv("OPENLUCKY_CONFIG")
	if configPath == "" {
		configPath = "openlucky.json"
	}

	addr := os.Getenv("OPENLUCKY_ADDR")
	if addr == "" {
		addr = "127.0.0.1:16601"
	}

	openLucky, err := app.New(app.Options{
		Addr:          addr,
		ConfigPath:    configPath,
		AdminUsername: os.Getenv("OPENLUCKY_ADMIN_USER"),
		AdminPassword: os.Getenv("OPENLUCKY_ADMIN_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("openlucky: %v", err)
	}

	log.Printf("openlucky listening on http://%s/lucky/", addr)
	openLucky.Spin()
}
