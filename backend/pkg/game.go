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

	// Basic game API

	// Choices handles the GET /choices request and returns the list of all choices.
	Choices(http.ResponseWriter, *http.Request)
	// Choice handles the GET /choice request and returns a random choice.
	Choice(http.ResponseWriter, *http.Request)
	// Play handles the POST /play request with users choice in the payload
	// and returns the game result and the choices made by the player and the computer.
	Play(http.ResponseWriter, *http.Request)

	// Scoreboard API

	// GetScores handles the GET /get_scores request and returns the list of all scores.
	GetScores(w http.ResponseWriter, r *http.Request)
	// ClearScores handles the POST /clear_scores request and clears the list of scores.
	ClearScores(w http.ResponseWriter, r *http.Request)

	// P2P API

	// CreateP2P handles the POST /create_p2p request and creates a new peer-to-peer game.
	CreateP2P(w http.ResponseWriter, r *http.Request)
	// ConnectP2P handles WebSocket connect request /connect_p2p with an existing peer-to-peer game.
	ConnectP2P(w http.ResponseWriter, r *http.Request)
	// FindP2PGame handles the GET /find_p2p request and returns the game status: full or not, if it is found.
	FindP2PGame(w http.ResponseWriter, r *http.Request)
}

// Storage is an interface that represents the storage of game results.
type Storage interface {
	// GetLastScores returns the list of last game results.
	GetLastScores() ([]types.Result, error)
	// SetLastScore adds a new game result to the list and removes the oldest one if the list is full.
	SetLastScore(types.Result) error
	// ClearScores clears the list of last game results.
	ClearScores() error
}

// The P2PGameFactory interface is for creating and managing peer-to-peer games. It has the following methods:
type P2PGameFactory interface {
	// CreateGame: This method creates a new peer-to-peer game with a given context and returns
	// the game object and an error if one occurred.
	CreateGame(ctx context.Context) (P2PGame, error)
	// StopGames: This method stops all games created by the factory.
	StopGames(ctx context.Context)
	// GetGame: This method retrieves a peer-to-peer game with a given ID. It returns
	// the game object and a boolean value indicating whether the game was found or not.
	GetGame(id types.GameID) (P2PGame, bool)
}

// P2PGame is an interface that represents a game played between two players.
type P2PGame interface {
	// GetID returns the unique identifier of the game.
	GetID() types.GameID
	// AddPlayer adds a player to the game with the given name. If no name is provided, the player will be called "Anonymous".
	// The name must contain only characters from the Latin alphabet and spaces, and must not be more than 20 characters long.
	// The function returns the side of the player (true = right side, false = left side) and a channel for receiving messages.
	//
	// The function will also send a signal to other player if one already joined
	//
	// Returns:
	// - side of the new player
	// - channel for current game state for the player and game results
	// - error if something is wrong
	AddPlayer(name string) (bool, chan types.Message, error)
	// RemovePlayer removes the player from the given side of the game.
	// The function will also send a signal to other player if one already joined
	RemovePlayer(rightSide bool)
	// Choice sets players choice on the given side of the game.
	// Sends the players the signal of current situation.
	// If both made choices, calculates result and sends it to the players.
	Choice(choice types.Choice, rightSide bool)
	// IsFull returns true if both players have joined the game.
	IsFull(ctx context.Context) bool
}
