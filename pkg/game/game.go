package game

import (
	"context"
	"fmt"

	"github.com/complynx/rpssl4bu/pkg"
)

type game struct {
	rng pkg.RandomProvider
}

var names = []string{
	"rock",
	"paper",
	"scissors",
	"lizard",
	"spock",
}

type pair struct {
	p1 int
	p2 int
}

var victoryMap = map[pair]bool{
	// against rock
	{0, 0}: false,
	{1, 0}: true,
	{2, 0}: false,
	{3, 0}: false,
	{4, 0}: true,

	// against paper
	{0, 1}: false,
	{1, 1}: false,
	{2, 1}: true,
	{3, 1}: true,
	{4, 1}: false,

	// against scissors
	{0, 2}: true,
	{1, 2}: false,
	{2, 2}: false,
	{3, 2}: false,
	{4, 2}: true,

	// against lizard
	{0, 3}: true,
	{1, 3}: false,
	{2, 3}: true,
	{3, 3}: false,
	{4, 3}: false,

	// against spock
	{0, 4}: false,
	{1, 4}: true,
	{2, 4}: false,
	{3, 4}: true,
	{4, 4}: false,
}

// returns tie and win results for player p1 against p2
func GameResult(p1 int, p2 int) (bool, bool) {
	if p1 == p2 {
		return true, false
	}
	return false, victoryMap[pair{
		p1: p1,
		p2: p2,
	}]
}

// returns tie and win results as string for player p1 against p2
func StringResult(p1 int, p2 int) string {
	tie, win := GameResult(p1, p2)
	if tie {
		return "tie"
	}
	if win {
		return "win"
	}
	return "lose"
}

func NewGame(rng pkg.RandomProvider) pkg.Game {
	return &game{
		rng: rng,
	}
}

func (g *game) Choices(ctx context.Context) ([]pkg.Choice, error) {
	ret := make([]pkg.Choice, len(names))
	for i, s := range names {
		ret[i] = pkg.Choice{
			ID:   i,
			Name: s,
		}
	}
	return ret, nil
}

func (g *game) Choice(ctx context.Context) (pkg.Choice, error) {
	num, err := g.rng.Rand(ctx)
	if err != nil {
		return pkg.Choice{}, fmt.Errorf("generate random number: %w", err)
	}

	numFactor := num % len(names)

	return pkg.Choice{
		ID:   numFactor,
		Name: names[numFactor],
	}, nil
}

func (g *game) Play(ctx context.Context, playerID int) (string, pkg.Choice, error) {

	computerChoice, err := g.Choice(ctx)
	if err != nil {
		return "", pkg.Choice{}, fmt.Errorf("get computer choice: %w", err)
	}

	return StringResult(playerID, computerChoice.ID), computerChoice, nil
}
