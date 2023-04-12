package main

import (
	"log"

	"github.com/Goboolean/data-control-server/internal"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
  }

	log.Fatal(internal.App())
}