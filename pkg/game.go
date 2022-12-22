package pkg

import (
	"context"
	"net/http"
)

//go:generate mockery --srcpkg github.com/complynx/rpssl4bu/pkg --all --with-expecter

type Server interface {
	Shutdown(ctx context.Context)
}

type Choice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RandomProvider interface {
	Rand(ctx context.Context) (int, error)
}

type Game interface {
	Choices(context.Context) ([]Choice, error)
	Choice(context.Context) (Choice, error)
	Play(context.Context, int) (string, Choice, error)
}

type GameAPI interface {
	Choices(http.ResponseWriter, *http.Request)
	Choice(http.ResponseWriter, *http.Request)
	Play(http.ResponseWriter, *http.Request)
}
