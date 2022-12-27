package gameapi

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/complynx/rpssl4bu/backend/pkg"
	"github.com/complynx/rpssl4bu/backend/pkg/types"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type gameAPI struct {
	game       pkg.Game
	p2pFactory pkg.P2PGameFactory
	log        *zap.Logger
	upgrader   websocket.Upgrader
	storage    pkg.Storage
}

func NewGameAPI(game pkg.Game, p2pFactory pkg.P2PGameFactory, storage pkg.Storage, log *zap.Logger) pkg.GameAPI {
	return &gameAPI{
		log:        log,
		game:       game,
		p2pFactory: p2pFactory,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		storage: storage,
	}
}

func (a *gameAPI) doRecover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		a.log.Error("Panic in api call", zap.Any("panic", r), zap.Any("stack_trace", buf[:stackSize]))
		httpCode(w, http.StatusInternalServerError)
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
	a.log.Info("sending Choices")

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

	if err == nil {
		a.log.Info("randomly chosen choice", zap.Any("computer_choice", choice))
	}

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

	if err == nil {
		if err := a.storage.SetLastScore(res); err != nil {
			a.log.Error("Failed to save last score",
				zap.Error(err),
			)
		}
		a.log.Info("game with computer",
			zap.Any("result", res),
			zap.Any("player_choice", req.Player),
			zap.Any("computer_choice", choice),
		)
	}

	a.marshalAndSend(playResult{
		Results:  res,
		Player:   req.Player.Int(),
		Computer: choice.Int(),
	}, err, w)
}

func (a *gameAPI) GetScores(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	scores, err := a.storage.GetLastScores()

	a.marshalAndSend(scores, err, w)
}

func (a *gameAPI) ClearScores(w http.ResponseWriter, r *http.Request) {
	defer a.doRecover(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	a.marshalAndSend(true, a.storage.ClearScores(), w)
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

type messageFromUser struct {
	Choice types.Choice `json:"choice"`
}

func sideString(isRight bool) string {
	if isRight {
		return "right"
	}
	return "left"
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
	log := a.log.With(
		zap.Any("game_id", id),
		zap.Any("name", name),
	)

	side, ch, err := game.AddPlayer(name)
	if err != nil {
		httpCode(w, http.StatusNotFound)
		return
	}
	defer game.RemovePlayer(side)

	log = log.With(zap.Any("side", sideString(side)))

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.sendErr(err, w, http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	go a.messageWriter(conn, side, log, ch)
	a.messageReader(conn, game, side, log)
}

type messageToUser struct {
	State types.Message `json:"state"`
	Side  string        `json:"side"`
}

func (a *gameAPI) messageWriter(conn *websocket.Conn, side bool, log *zap.Logger, ch <-chan types.Message) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			log.Error("Panic in messageWriter", zap.Any("panic", r), zap.Any("stack_trace", buf[:stackSize]))
			conn.Close()
		}
	}()

	for {
		msg, ok := <-ch
		if !ok {
			return
		}
		if side {
			msg.Result = msg.Result.Swap()
		}
		msgToUser := messageToUser{
			State: msg,
			Side:  sideString(side),
		}
		bytes, err := json.Marshal(msgToUser)
		if err != nil {
			log.Error("Error while marshalling message to json", zap.Error(err))
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Error("Error while sending message", zap.Error(err))
		}
	}
}

func (a *gameAPI) messageReader(conn *websocket.Conn, game pkg.P2PGame, side bool, log *zap.Logger) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error("Error while receiving message from websocket", zap.Error(err))
			return
		}
		var message messageFromUser
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Error("Error while unmarshalling message", zap.Error(err))
			continue
		}
		game.Choice(message.Choice, side)
	}
}
