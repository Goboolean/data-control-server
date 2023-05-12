package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/grpc"
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

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go server.Run(ctx, errCh, nil)

	select {
	case <- ctx.Done():
		if err := ctx.Err(); err != nil {
			log.Fatal(err)
		}
		cancel()
		wg.Wait()
	case <- signalCh:
		cancel()
		wg.Wait()
	}
}
