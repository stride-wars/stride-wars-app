package router

import (
	"net/http"

	apiroute "stride-wars-app/internal/api/apiconst"
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
	activityHandler *handler.ActivityHandler,
	activityService *service.ActivityService,
	hexLeaderboardHandler *handler.HexLeaderboardHandler,
	hexLeaderboardService *service.HexLeaderboardService,
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

	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()

	auth.HandleFunc(apiroute.Signup.String(), authHandler.SignUp).Methods("POST")
	auth.HandleFunc(apiroute.Signin.String(), authHandler.SignIn).Methods("POST")

	// User routes
	users := api.PathPrefix("/user").Subrouter()
	users.HandleFunc("", userHandler.GetUser).Methods("GET")
	users.HandleFunc(apiroute.UpdateUsername.String(), userHandler.UpdateUsername).Methods("PUT")

	// Activity routes
	activity := api.PathPrefix("/activity").Subrouter()
	activity.HandleFunc(apiroute.CreateActivity.String(), activityHandler.CreateActivity).Methods("POST")

	// Leaderboard routes
	leaderboard := api.PathPrefix("/leaderboard").Subrouter()
	leaderboard.HandleFunc(apiroute.GetLeaderboardByBBox.String(), hexLeaderboardHandler.GetAllLeaderboardsInsideBBox).Methods("GET")
	leaderboard.HandleFunc(apiroute.GetGlobalLeaderboard.String(), hexLeaderboardHandler.GetGlobalHexLeaderboard).Methods("GET")
}

// Handler returns the HTTP handler for the router
func (r *Router) Handler() http.Handler {
	return r.router
}
