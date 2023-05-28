package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)



func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}


	signalCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	/*
	 write constructor code here
	*/
	ctx, cancel := context.WithCancel(context.Background())

	/*
	 run handler with ctxx argument
	*/

	defer func () {
		// every fatel error will be catched here
		// call cancelFunc to cease all process
		if r := recover(); r != nil {
			cancel()		
		}
	}()

	select {
	case <- errCh:
	case <- ctx.Done():		
	}
}