package main

import (
	"log"
	"os"
	"os/signal"
	app "stride-wars-app/internal"
	"stride-wars-app/pkg/errors"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	app, err := initializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %s", err.Error())
	}
	app.Logger.Info("Application initialized.")

	if err := loadEnvs(); err != nil {
		app.Logger.Fatal("Failed to load environment variables", zap.Error(err))
	}

	if err := startApplication(app); err != nil {
		app.Logger.Fatal("Failed to start application", zap.Error(err))
	}
	app.Logger.Info("Application started.")

	waitForShutdown(app)
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

func startApplication(app *app.Application) error {
	if err := app.Start(); err != nil {
		return errors.WrapErr(err, "Failed to start application")
	}

	if err := app.StartHTTPServer(); err != nil {
		return errors.WrapErr(err, "Failed to start HTTP server")
	}

	app.Logger.Info("Application running...")
	return nil
}

func waitForShutdown(app *app.Application) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	app.Logger.Info("Shutting down...")
	if err := app.Stop(); err != nil {
		app.Logger.Fatal("Failed to gracefully shutdown application", zap.Error(err))
	}

	app.Logger.Info("Application stopped")
}
