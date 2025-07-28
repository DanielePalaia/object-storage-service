package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/yourusername/object-storage-service/docs" // Replace with your module path

	"github.com/yourusername/object-storage-service/domain"
)

type Server struct {
	router  *mux.Router
	storage domain.Storage
	port    string
}

// Package api implements HTTP handlers.
//
// @title Object Storage Service API
// @version 1.0
// @description API for storing, retrieving, and deleting objects.
// @host localhost:8080
// @BasePath /
//
// RegisterRoutes attaches HTTP handlers to the router
func RegisterRoutes(r *mux.Router, storage domain.Storage) {
	r.Use(loggingMiddleware)
	r.HandleFunc("/objects/{bucket}/{objectID}", putObjectHandler(storage)).Methods("PUT")
	r.HandleFunc("/objects/{bucket}/{objectID}", getObjectHandler(storage)).Methods("GET")
	r.HandleFunc("/objects/{bucket}/{objectID}", deleteObjectHandler(storage)).Methods("DELETE")
	r.HandleFunc("/health", HealthHandler).Methods("GET")
}

// NewServer creates a new server instance with storage and port config
func NewServer(storage domain.Storage, port string) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		storage: storage,
		port:    port,
	}

	// Register API routes
	RegisterRoutes(s.router, s.storage)

	// Register swagger UI route
	s.setupSwagger()

	return s
}

// Start runs the HTTP server
func (s *Server) Start() error {
	log.Printf("Server is running on port %s", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}

// setupSwagger adds Swagger UI handler on /docs/
func (s *Server) setupSwagger() {
	s.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
