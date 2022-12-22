package gameapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/complynx/rpssl4bu/pkg"
)

type gameAPI struct {
	game pkg.Game
}

func NewGameAPI(game pkg.Game) pkg.GameAPI {
	return &gameAPI{
		game: game,
	}
}

func doRecover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Panic", r)
	}
}

func sendErr(err error, w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	fmt.Println("Error: ", err)
}

func marshalAndSend(v any, err error, w http.ResponseWriter) {
	if err != nil {
		sendErr(err, w)
		return
	}

	resp, err := json.Marshal(v)
	if err != nil {
		sendErr(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *gameAPI) Choices(w http.ResponseWriter, r *http.Request) {
	defer doRecover(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	choices, err := a.game.Choices(r.Context())

	marshalAndSend(choices, err, w)
}

func (a *gameAPI) Choice(w http.ResponseWriter, r *http.Request) {
	defer doRecover(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	choice, err := a.game.Choice(r.Context())

	marshalAndSend(choice, err, w)
}

func (a *gameAPI) Play(w http.ResponseWriter, r *http.Request) {
	defer doRecover(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Player int `json:"player"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res, choice, err := a.game.Play(r.Context(), req.Player)

	marshalAndSend(struct {
		Results  string `json:"results"`
		Player   int    `json:"player"`
		Computer int    `json:"computer"`
	}{
		Results:  res,
		Player:   req.Player,
		Computer: choice.ID,
	}, err, w)
}
