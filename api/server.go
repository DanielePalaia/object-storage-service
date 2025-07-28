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
