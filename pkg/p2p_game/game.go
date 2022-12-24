package p2pgame

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/random"
	"github.com/complynx/rpssl4bu/pkg/types"
	"go.uber.org/zap"
)

const gameExistence = 2 * time.Hour

type player struct {
	Name string
	Chan chan types.Message
}

type p2pgame struct {
	ID      types.GameID
	factory *gameFactory
	log     *zap.Logger
	cancel  context.CancelFunc
	ctx     context.Context
	mu      sync.RWMutex

	p1 player
	p2 player
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
		factory: gf,
		log:     gf.log,

		p1: player{},
		p2: player{},
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

func (g *p2pgame) Start(ctx context.Context) error {
	g.log = g.log.With(zap.String("game_id", g.ID.String()))
	g.ctx, g.cancel = context.WithTimeout(context.Background(), gameExistence)
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

func (g *p2pgame) AddPlayer(name string) (bool, chan types.Message, error) {
	if !nameRe.MatchString(name) {
		return false, nil, ErrBadName
	}

	if name == "" {
		name = unnamed
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.p1.Name != "" {
		if g.p2.Name != "" {
			return false, nil, ErrGameIsFull
		}

		g.p2.Name = name
		g.p2.Chan = make(chan types.Message)
		return true, g.p2.Chan, nil
	}
	g.p1.Name = name
	g.p1.Chan = make(chan types.Message)
	return false, g.p1.Chan, nil
}

func (g *p2pgame) RemovePlayer(rightSide bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if rightSide {
		g.p2.Name = ""
		close(g.p2.Chan)
	}
	g.p1.Name = ""
	close(g.p1.Chan)
}

func (g *p2pgame) Choice(choice types.Choice, rightSide bool) {
}

func (g *p2pgame) run() {
	defer g.log.Info("p2p game finished")
	defer g.factory.removeGame(g.ID)
	defer g.cancel()

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			g.log.Error("Panic in game runner", zap.Any("panic", r), zap.Any("stack_trace", buf[:stackSize]))
		}
	}()

	g.log.Info("p2p game started")

	<-g.ctx.Done()
}
