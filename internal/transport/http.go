package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ntsiris/todo-app/api"
	"github.com/ntsiris/todo-app/internal/store"
)

type APIServer struct {
	*http.Server
	store store.Store
}

type APIServerConfig struct {
	Host string `envconfig:"API_SERVER_HOST" default:"localhost"`
	Port int    `envconfig:"API_SERVER_PORT" default:"8080"`
}

func NewAPIServer(config *APIServerConfig, store store.Store) *APIServer {
	return &APIServer{
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		},
		store: store,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	itemHandler := api.NewTodoHandler(s.store)
	itemHandler.RegisterRoutes(router)

	subRouter := http.NewServeMux()
	subRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	s.Handler = subRouter

	return s.ListenAndServe()
}

func (s *APIServer) Stop(ctx context.Context) error {
	if err := s.store.Close(); err != nil {
		return err
	}

	return s.Shutdown(ctx)
}
