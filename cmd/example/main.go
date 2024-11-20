package main

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/pkg/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	appserver, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	appserver.Start()

	if err := appserver.Stop(ctx); err != nil {
		return
	}

}
