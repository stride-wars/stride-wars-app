package internal

import (
	"context"
	"net/http"
	"os"
	"stride-wars-app/ent"
	"stride-wars-app/pkg/errors"

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
	Router         *http.Handler
}

func New() (*Application, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, errors.WrapErr(err, "Failed to initialize zap logger")
	}

	return &Application{Logger: logger}, nil
}

func (a *Application) Start() error {
	if err := a.initializeSupabaseClient(); err != nil {
		return err
	}

	if err := a.initializeEntClient(); err != nil {
		return err
	}

	if err := a.initializeRouter(); err != nil {
		return err
	}

	a.setupMuxRoutes()

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

func (a *Application) initializeEntClient() error {
	client, err := ent.Open("postgres", os.Getenv("SUPABASE_CONN_STRING"))
	if err != nil {
		return errors.WrapErr(err, "Failed to initialize Ent client")
	}

	// Ensure the schema is migrated
	// it's okay in development, but later on we might use ent/migrate
	if err := client.Schema.Create(context.Background()); err != nil {
		return errors.WrapErr(err, "Failed to create Ent schema")
	}

	a.EntClient = client
	return nil
}

func (a *Application) initializeRouter() error {
	m := mux.NewRouter()
	a.Router = applyCors(m)
	return nil
}

func (a *Application) setupMuxRoutes() {
	// Will require mux route setup later on similar to this:
	// a.Router.HandleFunc("/hex", a.CreateHexHandler).Methods("POST")
}

func (a *Application) StartHTTPServer() error {
	a.Logger.Info("Starting HTTP server")
	if err := http.ListenAndServe(":8080", *a.Router); err != nil {
		return errors.WrapErr(err, "Failed when starting HTTP server")
	}
	return nil
}

func (a *Application) Stop() error {
	if err := a.EntClient.Close(); err != nil {
		return errors.WrapErr(err, "Failed to close Ent client connection")
	}

	a.Logger.Info("Database connection closed.")
	a.Logger.Info("Application stopped gracefully.")
	return nil
}

// applyCors sets up CORS middleware for the router.
// might require different CORS for different origins?
func applyCors(router *mux.Router) *http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO restrict this for specific domains
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	r := c.Handler(router)
	return &r
}
