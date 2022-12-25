package pkg

import (
	"context"
	"net/http"

	"github.com/complynx/rpssl4bu/backend/pkg/types"
)

//go:generate mockery --srcpkg github.com/complynx/rpssl4bu/backend/pkg --all --with-expecter

// Server is an interface that represents a server that can be shut down.
type Server interface {
	// Shutdown shuts down the server.
	// The provided context is used to cancel any ongoing operations.
	Shutdown(ctx context.Context)
}

// RandomProvider is an interface that represents a provider of random numbers.
type RandomProvider interface {
	// Rand returns a random number from 0 to 99
	// The provided context is used to cancel the request if it takes too long.
	Rand(ctx context.Context) (int, error)
}

// Game is an interface that represents a game.
type Game interface {
	// Choices returns the list of all choices.
	Choices(context.Context) ([]types.Choice, error)
	// Choice returns a random choice.
	Choice(context.Context) (types.Choice, error)
	// Play runs the game based on users choice and returns the game result and
	// the choice made by the the computer.
	Play(context.Context, types.Choice) (types.Result, types.Choice, error)
}

// GameAPI is an interface that represents the API of the game.
type GameAPI interface {
	// Choices handles the GET /choices request and returns the list of all choices.
	Choices(http.ResponseWriter, *http.Request)
	// Choice handles the GET /choice request and returns a random choice.
	Choice(http.ResponseWriter, *http.Request)
	// Play handles the POST /play request with users choice in the payload
	// and returns the game result and the choices made by the player and the computer.
	Play(http.ResponseWriter, *http.Request)
}

type Storage interface {
	GetLastScores() ([]types.Result, error)
	SetLastScore(types.Result) error
	ClearScores() error
}

type P2PGameFactory interface {
	CreateGame(ctx context.Context) (P2PGame, error)
	StopGames(ctx context.Context)
	GetGame(id types.GameID) (P2PGame, bool)
}

type P2PGame interface {
	GetID() types.GameID
	// adds player with given name, name is checked to be of latin alphabet (+spaces) and not more than
	// 20 symbols. If no name provided, user is called Anonymous.
	// returns the side of the player true = right
	// and channel for messages
	AddPlayer(name string) (bool, chan types.Message, error)
	RemovePlayer(rightSide bool)
	Choice(choice types.Choice, rightSide bool)
}
