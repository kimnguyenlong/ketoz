package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/config"
	"github.com/kimnguyenlong/ketoz/internal/handler"
	"github.com/kimnguyenlong/ketoz/internal/repository"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
)

func main() {
	// Run Keto migrations
	slog.Info("Running Keto migrations...")
	if err := runKetoMigrations(); err != nil {
		slog.Error("Failed to run Keto migrations", "error", err)
		return
	}
	slog.Info("Keto migrations completed successfully")

	// Start Keto server
	slog.Info("Starting Keto server...")
	ketoCmd, err := runKeto()
	defer func() {
		if err := ketoCmd.Process.Signal(syscall.SIGTERM); err != nil {
			slog.Error("Failed to kill Keto process", "error", err)
		}
		ketoCmd.Wait()
		slog.Info("Keto server stopped gracefully")
	}()
	if err != nil {
		slog.Error("Failed to start Keto", "error", err)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	keto, err := keto.NewKeto()
	if err != nil {
		slog.Error("Failed to init Keto client", "error", err)
		return
	}

	app := fiber.New()
	api := app.Group("/api")

	idRepo := repository.NewIdentity(keto)
	roleRepo := repository.NewRole(keto)
	rscRepo := repository.NewResource(keto)
	pmRepo := repository.NewPermission(keto)

	idHandler := handler.NewIdentity(idRepo)
	roleHandler := handler.NewRole(roleRepo)
	rscHandler := handler.NewResource(rscRepo)
	pmHandler := handler.NewPermission(pmRepo)

	idHandler.RegisterRoutes(api)
	roleHandler.RegisterRoutes(api)
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

func runKetoMigrations() error {
	ketoCmd := exec.Command("keto", "migrate", "up", "-c", "/home/ory/config.yml", "-y")
	ketoCmd.Stdout = os.Stdout
	ketoCmd.Stderr = os.Stderr
	return ketoCmd.Run()
}

func runKeto() (*exec.Cmd, error) {
	ketoCmd := exec.Command("keto", "serve", "-c", "/home/ory/config.yml")
	ketoCmd.Stdout = os.Stdout
	ketoCmd.Stderr = os.Stderr
	err := ketoCmd.Start()
	return ketoCmd, err
}
