package main

import (
	"log"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/grpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	log.Fatal(server.Run())
}
