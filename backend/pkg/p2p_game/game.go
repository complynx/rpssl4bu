package p2pgame

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/complynx/rpssl4bu/backend/pkg"
	"github.com/complynx/rpssl4bu/backend/pkg/game"
	"github.com/complynx/rpssl4bu/backend/pkg/random"
	"github.com/complynx/rpssl4bu/backend/pkg/types"
	"go.uber.org/zap"
)

const gameExistence = 2 * time.Hour

type player struct {
	Name   string
	Choice types.Choice
	Chan   chan types.Message
}

type p2pgame struct {
	ID          types.GameID
	factory     *gameFactory
	log         *zap.Logger
	cancel      context.CancelFunc
	ctx         context.Context
	mu          sync.RWMutex
	pingChannel chan struct{}

	left  player
	right player
}

type gameFactory struct {
	rng   pkg.RandomProvider
	games map[types.GameID]*p2pgame
	mu    sync.RWMutex
	log   *zap.Logger
}

func NewGameFactory(rng pkg.RandomProvider, log *zap.Logger) pkg.P2PGameFactory {
	return &gameFactory{
		rng:   rng,
		games: make(map[types.GameID]*p2pgame),
		log:   log,
	}
}

func (gf *gameFactory) setGameIfNotExist(g *p2pgame) bool {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	_, exists := gf.games[g.ID]
	if exists {
		return false
	}

	gf.games[g.ID] = g
	return true
}

func (gf *gameFactory) StopGames(ctx context.Context) {
	gf.mu.RLock()
	defer gf.mu.RUnlock()

	for _, g := range gf.games {
		g.cancel()
	}
}

func (gf *gameFactory) removeGame(id types.GameID) {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	delete(gf.games, id)
}

func (gf *gameFactory) CreateGame(ctx context.Context) (pkg.P2PGame, error) {
	var id types.GameID
	var err error
	game := &p2pgame{
		factory:     gf,
		log:         gf.log,
		pingChannel: make(chan struct{}),

		left:  player{},
		right: player{},
	}
	for {
		id, err = random.RandomID(ctx, gf.rng)
		if err != nil {
			return nil, fmt.Errorf("generating ID: %w", err)
		}
		game.ID = id
		if gf.setGameIfNotExist(game) {
			break
		}
	}
	err = game.Start(ctx)
	if err != nil {
		gf.removeGame(game.ID)
		return nil, fmt.Errorf("starting game: %w", err)
	}
	return game, nil
}

func (gf *gameFactory) GetGame(id types.GameID) (pkg.P2PGame, bool) {
	gf.mu.RLock()
	defer gf.mu.RUnlock()

	ret, ok := gf.games[id]
	return ret, ok
}

func (g *p2pgame) IsFull(ctx context.Context) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.left.Name != "" && g.right.Name != ""
}

func (g *p2pgame) Start(ctx context.Context) error {
	g.log = g.log.With(zap.String("game_id", g.ID.String()))
	g.ctx, g.cancel = context.WithCancel(context.Background())
	go g.run()
	return nil
}

func (g *p2pgame) GetID() types.GameID {
	return g.ID
}

var ErrBadName = fmt.Errorf("bad name")
var ErrGameIsFull = fmt.Errorf("game is full")
var nameRe = regexp.MustCompile(`^[a-zA-Z ]{0,20}$`)

const unnamed = "Anonymous"

func (g *p2pgame) AddPlayer(name string) (rightSide bool, ch chan types.Message, err error) {
	if !nameRe.MatchString(name) {
		return false, nil, ErrBadName
	}

	if name == "" {
		name = unnamed
	}

	defer g.sendState(types.Unknown)
	go g.ping()

	defer func() {
		if err == nil {
			g.log.Info("player added",
				zap.Bool("side", rightSide),
				zap.Any("player1", g.left.Name),
				zap.Any("player2", g.right.Name),
			)
		}
	}()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.left.Name != "" {
		if g.right.Name != "" {
			return false, nil, ErrGameIsFull
		}

		g.right.Name = name
		g.right.Chan = make(chan types.Message)
		return true, g.right.Chan, nil
	}
	g.left.Name = name
	g.left.Chan = make(chan types.Message)
	return false, g.left.Chan, nil
}

func (g *p2pgame) sendState(result types.Result) {
	g.log.Info("sending state",
		zap.Any("result", result),
		zap.Any("player1", g.left.Name),
		zap.Any("player2", g.right.Name),
		zap.Any("player1_choice", g.left.Choice),
		zap.Any("player2_choice", g.right.Choice),
	)

	if g.left.Name != "" {
		go func(n1, n2 string, c1, c2 types.Choice, result types.Result) {
			if c1 == types.Undefined {
				c2 = types.Undefined
			}
			g.left.Chan <- types.Message{
				LeftPlayerName:    n1,
				RightPlayerName:   n2,
				Result:            result,
				LeftPlayerChoice:  c1,
				RightPlayerChoice: c2,
			}
			g.log.Info("sent state to p1",
				zap.Any("result", result),
				zap.Any("player1", n1),
				zap.Any("player2", n2),
				zap.Any("player1_choice", c1),
				zap.Any("player2_choice", c2),
			)
		}(g.left.Name, g.right.Name, g.left.Choice, g.right.Choice, result)
	}
	if g.right.Name != "" {
		go func(n1, n2 string, c1, c2 types.Choice, result types.Result) {
			if c2 == types.Undefined {
				c1 = types.Undefined
			}
			g.right.Chan <- types.Message{
				LeftPlayerName:    n1,
				RightPlayerName:   n2,
				Result:            result.Swap(),
				LeftPlayerChoice:  c1,
				RightPlayerChoice: c2,
			}
			g.log.Info("sent state to p2",
				zap.Any("result", result.Swap()),
				zap.Any("player1", n1),
				zap.Any("player2", n2),
				zap.Any("player1_choice", c1),
				zap.Any("player2_choice", c2),
			)
		}(g.left.Name, g.right.Name, g.left.Choice, g.right.Choice, result)
	}
}

func (g *p2pgame) ping() {
	select {
	case g.pingChannel <- struct{}{}:
	default:
	}
}

func (g *p2pgame) RemovePlayer(rightSide bool) {
	defer g.sendState(types.Unknown)

	defer func() {
		g.log.Info("player removed",
			zap.Bool("side", rightSide),
			zap.Any("player1", g.left.Name),
			zap.Any("player2", g.right.Name),
		)
	}()

	g.mu.Lock()
	defer g.mu.Unlock()

	if rightSide {
		g.right.Name = ""
		close(g.right.Chan)
	} else {
		g.left.Name = ""
		close(g.left.Chan)
	}
}

func (g *p2pgame) Choice(choice types.Choice, rightSide bool) {
	go g.ping()

	g.log.Info("User choice", zap.Bool("side", rightSide), zap.Any("choice", choice))
	res := types.Unknown

	g.mu.Lock()
	defer g.mu.Unlock()

	if rightSide {
		g.right.Choice = choice
	} else {
		g.left.Choice = choice
	}
	if g.left.Choice != types.Undefined && g.right.Choice != types.Undefined {
		res = game.GameResult(g.left.Choice, g.right.Choice)
	}
	g.sendState(res)
	if res != types.Unknown {
		g.left.Choice = types.Undefined
		g.right.Choice = types.Undefined
	}
}

func (g *p2pgame) run() {
	defer g.log.Info("p2p game finished")
	defer g.factory.removeGame(g.ID)
	defer g.cancel()
	defer close(g.pingChannel)

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			g.log.Error("Panic in game runner", zap.Any("panic", r), zap.Any("stack_trace", buf[:stackSize]))
		}
	}()

	g.log.Info("p2p game started")

	timer := time.NewTimer(gameExistence)

	for {
		select {
		case <-g.pingChannel:
			timer.Reset(gameExistence)
		case <-timer.C:
			return
		case <-g.ctx.Done():
			return
		}
	}
}
