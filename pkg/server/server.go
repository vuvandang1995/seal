package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vuvandang1995/seal/pkg/config"
)

// Server structure
type Server struct {
	router chi.Router
	config *config.Config
}

func New() *Server {
	server := &Server{}
	server.router = chi.NewRouter()
	return server
}

func (s *Server) Setup() {
	s.config = config.Read()
	s.routes()
}

// ServeHTTP provide method to serve HTTP Handle
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
