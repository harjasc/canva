package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
}

type Options struct {
	host string
	port int
}

func New(opts Options) *Server {
	address := fmt.Sprintf("%s:%d", opts.host, opts.port)
	mux := chi.NewMux()
	return &Server{
		address: address,
		mux:     mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}

}

// Starts the server, setting up routes and listens for HTTP requests
func (s *Server) Start() error {
	s.setupRoutes()

	fmt.Println("Starting on", s.address)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Gracefully stops the server within the timeout
func (s *Server) Stop() error {
	fmt.Println("Stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}
	return nil
}
