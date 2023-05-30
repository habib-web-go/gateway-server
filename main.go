package main

import (
	"log"

	redis "github.com/habib-web-go/gateway-server/redis"
	routes "github.com/habib-web-go/gateway-server/routes"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	redis.RunRedis()
	routes.Run()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
