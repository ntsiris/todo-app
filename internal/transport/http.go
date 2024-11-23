package transport

import (
	"context"
	"github.com/ntsiris/todo-app/api"
	"github.com/ntsiris/todo-app/internal/store"
	"net/http"
)

type APIServer struct {
	*http.Server
	store store.Store
}

func NewAPIServer(address string, store store.Store) *APIServer {
	return &APIServer{
		Server: &http.Server{
			Addr: address,
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
