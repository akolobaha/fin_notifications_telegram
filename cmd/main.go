package main

import (
	"context"
	"fin_notifications_telegram/cmd/commands"
	"fin_notifications_telegram/internal/config"
	"fin_notifications_telegram/internal/log"
	"fin_notifications_telegram/internal/monitoring"
	"os"
	"os/signal"
	"syscall"
)

const defaultEnvFilePath = ".env"

func init() {
	monitoring.RegisterPrometheus()
}

func main() {
	cfg, err := config.Parse(defaultEnvFilePath)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	monitoring.RunPrometheusServer(cfg.GetPrometheusURL())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	err = commands.ReadFromQueue(ctx, cfg)
	if err != nil {
		log.Error("Error reading from queue", err)
		return
	}

}
