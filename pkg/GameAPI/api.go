package gameapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/complynx/rpssl4bu/pkg"
	"go.uber.org/zap"
)

type gameAPI struct {
	game pkg.Game
	log  *zap.Logger
}

func NewGameAPI(game pkg.Game, log *zap.Logger) pkg.GameAPI {
	return &gameAPI{
		log:  log,
		game: game,
	}
}

func (a *gameAPI) doRecover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		httpCode(w, http.StatusInternalServerError)
		a.log.Panic(fmt.Sprintf("%s", r), zap.Any("stack_trace", buf[:stackSize]))
	}
}

func httpCode(w http.ResponseWriter, errCode int) {
	http.Error(w, http.StatusText(errCode), errCode)
}

func (a *gameAPI) sendErr(err error, w http.ResponseWriter, errCode int) {
	httpCode(w, errCode)
	a.log.Error("Error during request processing", zap.Error(err))
}

func (a *gameAPI) marshalAndSend(v any, err error, w http.ResponseWriter) {
	if err != nil {
		a.sendErr(err, w, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(v)
	if err != nil {
		a.sendErr(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *gameAPI) Choices(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	choices, err := a.game.Choices(r.Context())

	a.marshalAndSend(choices, err, w)
}

func (a *gameAPI) Choice(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	choice, err := a.game.Choice(r.Context())

	a.marshalAndSend(choice, err, w)
}

type playResult struct {
	Results  string `json:"results"`
	Player   int    `json:"player"`
	Computer int    `json:"computer"`
}

func (a *gameAPI) Play(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

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

	a.marshalAndSend(playResult{
		Results:  res,
		Player:   req.Player,
		Computer: choice.ID,
	}, err, w)
}
