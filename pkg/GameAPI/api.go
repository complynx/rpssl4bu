package gameapi

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/types"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type gameAPI struct {
	game       pkg.Game
	p2pFactory pkg.P2PGameFactory
	log        *zap.Logger
	upgrader   websocket.Upgrader
}

func NewGameAPI(game pkg.Game, p2pFactory pkg.P2PGameFactory, log *zap.Logger) pkg.GameAPI {
	return &gameAPI{
		log:        log,
		game:       game,
		p2pFactory: p2pFactory,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (a *gameAPI) doRecover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		httpCode(w, http.StatusInternalServerError)
		a.log.Error("Panic in api call", zap.Any("panic", r), zap.Any("stack_trace", buf[:stackSize]))
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
	Results  types.Result `json:"results"`
	Player   int          `json:"player"`
	Computer int          `json:"computer"`
}

func (a *gameAPI) Play(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Player types.Choice `json:"player"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res, choice, err := a.game.Play(r.Context(), req.Player)

	a.marshalAndSend(playResult{
		Results:  res,
		Player:   req.Player.Int(),
		Computer: choice.Int(),
	}, err, w)
}

func (a *gameAPI) CreateP2P(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	game, err := a.p2pFactory.CreateGame(r.Context())
	if err != nil {
		a.sendErr(err, w, http.StatusInternalServerError)
		return
	}

	a.marshalAndSend(game.GetID(), err, w)
}

func (a *gameAPI) ConnectP2P(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	id, err := types.GameIDFromString(r.URL.Query().Get("g"))
	if err != nil {
		httpCode(w, http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")

	game, found := a.p2pFactory.GetGame(id)
	if !found {
		httpCode(w, http.StatusNotFound)
		return
	}

	side, ch, err := game.AddPlayer(name)
	if err != nil {
		httpCode(w, http.StatusNotFound)
		return
	}
	defer game.RemovePlayer(side)

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.sendErr(err, w, http.StatusInternalServerError)
		return
	}
	defer conn.Close()
}
