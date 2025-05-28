package router

import (
	"net/http"

	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Router represents the API router
type Router struct {
	router *mux.Router
	logger *zap.Logger
}

func New(logger *zap.Logger) *Router {
	return &Router{
		router: mux.NewRouter(),
		logger: logger,
	}
}

// Setup configures all routes and middleware
func (r *Router) Setup(
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	userHandler *handler.UserHandler,
	userService *service.UserService,
) {
	// CORS must be first to handle preflight requests
	r.router.Use(middleware.CORS())

	// Debug middleware to log all requests
	r.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			r.logger.Info("Incoming request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("host", req.Host),
			)
			next.ServeHTTP(w, req)
		})
	})

	r.router.Use(middleware.Logger(r.logger))
	r.router.Use(middleware.ErrorHandler(r.logger))
	r.router.Use(middleware.ParseJSON)

	// Handle OPTIONS requests globally
	r.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Set up NotFound and MethodNotAllowed handlers
	r.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.logger.Warn("Route not found",
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
		)
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	r.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.logger.Warn("Method not allowed",
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
		)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// API v1 routes
	api := r.router.PathPrefix("/api/v1").Subrouter()

	// Test route
	api.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Router is working!"))
	}).Methods("GET")

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	auth.HandleFunc("/signin", authHandler.SignIn).Methods("POST")

	users := api.PathPrefix("/user").Subrouter()
	users.HandleFunc("/{username}", userHandler.GetUserByUsername).Methods("GET")
	users.HandleFunc("/{username}", userHandler.UpdateUsername).Methods("PUT")
	users.HandleFunc("/{user-id}",userHandler.GetUserByID).Methods("GET")
}

// Handler returns the HTTP handler for the router
func (r *Router) Handler() http.Handler {
	return r.router
}
