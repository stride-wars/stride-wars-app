package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"stride-wars-app/ent"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
	"stride-wars-app/pkg/errors"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type Application struct {
	Logger         *zap.Logger
	SupabaseClient *supabase.Client
	EntClient      *ent.Client
	Router         http.Handler // use interface, NOT pointer

	Services     *service.Services
	Handlers     *handler.Handlers
	Repositories *repository.Repositories
}

func New() (*Application, error) {
	logger, err := zap.NewProduction()
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

	err := a.initializeRouter()

	return err
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
	m := mux.NewRouter()

	m.HandleFunc("/api/auth/signup", a.Handlers.AuthHandler.SignUp).Methods("POST")
	m.HandleFunc("/api/auth/signin", a.Handlers.AuthHandler.SignIn).Methods("POST")

	a.Router = applyCors(m)
	return nil
}

func (a *Application) StartHTTPServer() error {
	a.Logger.Info("Starting HTTP server")

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.Router,
	}
	errChan := make(chan error, 1)

	go func() {
		a.Logger.Info("HTTP server listening on :8080")
		errChan <- server.ListenAndServe()
	}()

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

func applyCors(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	return c.Handler(handler)
}
