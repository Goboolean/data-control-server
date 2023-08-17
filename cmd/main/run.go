package main

import (
	"log"

	"github.com/Goboolean/fetch-server/cmd/compose"
	_ "github.com/Goboolean/fetch-server/internal/util/logger"
	"github.com/joho/godotenv"
)



func main() {
	godotenv.Load()
	log.Fatal(compose.Released())
}