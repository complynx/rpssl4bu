package server

import (
	"context"
	"net/http"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type server struct {
	srv *http.Server
	log *zap.Logger
}

func StartHTTPServer(listen string, api pkg.GameAPI, log *zap.Logger) pkg.Server {
	mux := setupRouter(api)
	srv := &server{
		srv: &http.Server{
			Addr:    listen,
			Handler: mux,
		},
		log: log,
	}

	go func() {
		if err := srv.srv.ListenAndServe(); err != nil {
			log.Error("Server error", zap.Error(err))
		}
	}()

	log.Info("Server started", zap.Any("address", listen))

	return srv
}

func setupRouter(api pkg.GameAPI) *chi.Mux {
	httpRouter := chi.NewMux()

	httpRouter.Use(
		WithAccessControlAllowOrigin(),
	)

	httpRouter.HandleFunc("/choices", api.Choices)
	httpRouter.HandleFunc("/choice", api.Choice)
	httpRouter.HandleFunc("/play", api.Play)

	return httpRouter
}

func (srv *server) Shutdown(ctx context.Context) {
	srv.srv.Shutdown(ctx)

	srv.log.Info("Server stopped")
}
