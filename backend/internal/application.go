package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"stride-wars-app/ent"
	"stride-wars-app/internal/api/router"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
	"stride-wars-app/pkg/errors"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Logger         *zap.Logger
	SupabaseClient *supabase.Client
	EntClient      *ent.Client
	Router         http.Handler

	Services     *service.Services
	Handlers     *handler.Handlers
	Repositories *repository.Repositories
}

func New() (*Application, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, errors.WrapErr(err, "Failed to initialize zap logger")
	}

	return &Application{Logger: logger}, nil
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.initializeSupabaseClient(); err != nil {
		return err
	}

	if err := a.initializeEntClient(ctx); err != nil {
		return err
	}

	a.Repositories = repository.Provide(a.EntClient)
	a.Services = service.Provide(a.Repositories, a.SupabaseClient, a.Logger)
	a.Handlers = handler.Provide(a.Services, a.Logger)

	if err := a.initializeRouter(); err != nil {
		return err
	}

	return nil
}

func (a *Application) initializeSupabaseClient() error {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_PROJECT_URL"), os.Getenv("SUPABASE_API_KEY"), &supabase.ClientOptions{})
	if err != nil {
		return errors.WrapErr(err, "Failed to initialize Supabase client")
	}
	a.SupabaseClient = client
	return nil
}

func (a *Application) initializeEntClient(ctx context.Context) error {
	client, err := ent.Open("postgres", os.Getenv("SUPABASE_CONN_STRING"))
	if err != nil {
		return errors.WrapErr(err, "Failed to initialize Ent client")
	}

	if err := client.Schema.Create(ctx); err != nil {
		return errors.WrapErr(err, "Failed to create Ent schema")
	}

	a.EntClient = client
	return nil
}

func (a *Application) initializeRouter() error {
	router := router.New(a.Logger)
	router.Setup(a.Handlers.AuthHandler, a.Services.AuthService)
	a.Router = router.Handler()
	return nil
}

func (a *Application) StartHTTPServer() error {
	a.Logger.Info("Starting HTTP server")

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.Router,
	}
	errChan := make(chan error, 1)
	readyChan := make(chan struct{})

	go func() {
		a.Logger.Info("HTTP server listening on :8080")
		close(readyChan) // Signal that the server is about to start
		errChan <- server.ListenAndServe()
	}()

	// Wait for the server to be ready
	<-readyChan

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		a.Logger.Info("Received signal", zap.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return errors.WrapErr(err, "Failed to gracefully shutdown HTTP server")
		}
		a.Logger.Info("HTTP server shut down gracefully")
		return nil

	case err := <-errChan:
		return errors.WrapErr(err, "HTTP server encountered error")
	}
}

func (a *Application) Stop() error {
	if err := a.EntClient.Close(); err != nil {
		return errors.WrapErr(err, "Failed to close Ent client connection")
	}

	a.Logger.Info("Database connection closed.")
	a.Logger.Info("Application stopped gracefully.")
	return nil
}
