package main

import (
	"context"
	"github.com/khostya/pvz/internal/app"
	"github.com/khostya/pvz/internal/config"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = app.StartApp(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
