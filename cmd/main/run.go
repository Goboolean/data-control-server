package main

import (
	"context"
	"log"

	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	/*
	 write constructor code here
	*/
	ctx, cancel := context.WithCancel(context.Background())

	/*
	 run handler giving ctx argument

	 prometheus / grpc
	 buycycle / polygon
	*/

	prometheus.Run(ctx)

	defer func() {
		// every fatel error will be catched here
		// call cancelFunc to cease all process
		if err := recover(); err != nil {
			log.Fatal(err)
		}

		cancel()
	}()

	<-ctx.Done()
}
