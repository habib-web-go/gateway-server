package main

import (
	"log"

	routes "github.com/habib-web-go/gateway-server/routes"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	routes.Run()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
