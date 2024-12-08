package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Utro-tvar/medods-test/internal/service"
	"github.com/Utro-tvar/medods-test/internal/storage/postgres"
	"github.com/Utro-tvar/medods-test/internal/transport/rest"
)

func main() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("Starting token service")

	storage, err := postgres.New()
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("storage initialized")

	service := service.New(log, storage)

	log.Info("service initialized")

	rest := rest.New(log, service)

	go rest.MustRun("localhost:8080")

	log.Info("service start listening at localhost:8080")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	log.Info("stopping server")

	storage.Close()

	log.Info("server stopped")
}
