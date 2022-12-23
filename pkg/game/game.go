package game

import (
	"context"
	"fmt"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/types"
)

type game struct {
	rng pkg.RandomProvider
}

type pair struct {
	p1 types.Choice
	p2 types.Choice
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
func GameResult(p1, p2 types.Choice) types.Result {
	if p1 == p2 {
		return types.Tie
	}
	if victoryMap[pair{
		p1: p1,
		p2: p2,
	}] {
		return types.Win
	}
	return types.Lose
}

func NewGame(rng pkg.RandomProvider) pkg.Game {
	return &game{
		rng: rng,
	}
}

func (g *game) Choices(ctx context.Context) ([]types.Choice, error) {
	ret := make([]types.Choice, 0, 5)
	for i := 0; i < 5; i++ {
		ret = append(ret, types.IntToChoice(i))
	}
	return ret, nil
}

func (g *game) Choice(ctx context.Context) (types.Choice, error) {
	num, err := g.rng.Rand(ctx)
	if err != nil {
		return types.Lizard, fmt.Errorf("generate random number: %w", err)
	}

	return types.IntToChoiceErr(num % 5)
}

func (g *game) Play(ctx context.Context, player types.Choice) (types.Result, types.Choice, error) {

	computerChoice, err := g.Choice(ctx)
	if err != nil {
		return types.Tie, types.Spock, fmt.Errorf("get computer choice: %w", err)
	}

	return GameResult(player, computerChoice), computerChoice, nil
}
