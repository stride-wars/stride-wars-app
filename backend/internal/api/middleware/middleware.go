package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// contextKey is a type for context keys
type contextKey string

// JSONBodyKey is the key used to store JSON body in context
const JSONBodyKey contextKey = "json_body"

// Response represents a standard API response
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ParseJSON middleware parses JSON request body
func ParseJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "DELETE" {
			next.ServeHTTP(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Failed to read request body")
			return
		}
		if err := r.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}

		if len(body) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		var data any
		if err := json.Unmarshal(body, &data); err != nil {
			WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
			return
		}

		// Store parsed JSON in context for handlers to use
		ctx := context.WithValue(r.Context(), JSONBodyKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetJSONBody retrieves parsed JSON from context
func GetJSONBody(r *http.Request) (any, bool) {
	data := r.Context().Value(JSONBodyKey)
	if data == nil {
		return nil, false
	}
	return data, true
}

// Logger middleware logs request details
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			// Log the request details
			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.String("duration", time.Since(start).String()),
				zap.String("ip", r.RemoteAddr),
			)
		})
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for all responses
			//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
			w.Header().Set("Access-Control-Allow-Origin", "http://192.168.22.230:8081")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// For non-preflight requests, continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// ErrorHandler middleware handles panics and errors
func ErrorHandler(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
					)

					WriteError(w, http.StatusInternalServerError, "Internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// WriteJSON writes a JSON response to the client
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// WriteError writes an error response to the client
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
