package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/go-chi/chi/v5"
)

type server struct {
	srv *http.Server
}

func StartHTTPServer(port int, api pkg.GameAPI) pkg.Server {
	mux := setupRouter(api)
	srv := &server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	go func() {
		if err := srv.srv.ListenAndServe(); err != nil {
			fmt.Println("Server error: %w", err)
		}
	}()

	return srv
}

func setupRouter(api pkg.GameAPI) *chi.Mux {
	httpRouter := chi.NewMux()

	httpRouter.HandleFunc("/choices", api.Choices)
	httpRouter.HandleFunc("/choice", api.Choice)
	httpRouter.HandleFunc("/play", api.Play)

	return httpRouter
}

func (srv *server) Shutdown(ctx context.Context) {
	srv.srv.Shutdown(ctx)
}
