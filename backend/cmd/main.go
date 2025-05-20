package main

import (
	"context"
	"log"
	"os/signal"
	app "stride-wars-app/internal"
	"stride-wars-app/pkg/errors"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := loadEnvs(); err != nil {
		log.Printf("Warning: %v\n", err)
	}

	app, err := initializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %s", err.Error())
	}
	app.Logger.Info("Application initialized.")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := startApplication(ctx, app); err != nil {
		app.Logger.Fatal("Failed to start application", zap.Error(err))
	}

	app.Logger.Info("Application running...")

	<-ctx.Done()

	app.Logger.Info("Shutting down...")
	if err := app.Stop(); err != nil {
		app.Logger.Fatal("Failed to gracefully shutdown application", zap.Error(err))
	}

	app.Logger.Info("Application stopped")
}

func initializeApplication() (*app.Application, error) {
	app, err := app.New()
	if err != nil {
		return nil, errors.WrapErr(err, "Failed to initialize application")
	}
	return app, nil
}

func loadEnvs() error {
	if err := godotenv.Load(); err != nil {
		return errors.WrapErr(err, "Failed to load .env file")
	}
	return nil
}

func startApplication(ctx context.Context, app *app.Application) error {
	if err := app.Start(ctx); err != nil {
		return errors.WrapErr(err, "Failed to start application")
	}

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- app.StartHTTPServer()
	}()

	// Wait for either context cancellation or server error
	select {
	case err := <-serverErrChan:
		return errors.WrapErr(err, "HTTP server error")
	case <-ctx.Done():
		return nil
	}
}
