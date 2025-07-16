package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/config"
	"github.com/kimnguyenlong/ketoz/internal/handler"
	"github.com/kimnguyenlong/ketoz/internal/repository"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	keto, err := keto.NewKeto(cfg.Keto)
	if err != nil {
		slog.Error("Failed to init Keto client", "error", err)
		return
	}
	defer keto.Close()

	app := fiber.New()
	api := app.Group("/api")

	idRepo := repository.NewIdentity(keto)
	rscRepo := repository.NewResource(keto)
	pmRepo := repository.NewPermission(keto)

	idHandler := handler.NewIdentity(idRepo)
	rscHandler := handler.NewResource(rscRepo)
	pmHandler := handler.NewPermission(pmRepo)

	idHandler.RegisterRoutes(api)
	rscHandler.RegisterRoutes(api)
	pmHandler.RegisterRoutes(api)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-stop
		slog.Info("Shutting down server gracefully...")
		app.Shutdown()
	}()

	slog.Info("Starting server", "host", cfg.Service.Host, "port", cfg.Service.Port)
	defer slog.Info("Server stopped gracefully")
	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)); err != nil {
		slog.Error("Serve HTTP requests failed", "error", err)
		return
	}
}
