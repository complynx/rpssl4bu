package game

import (
	"context"
	"fmt"

	"github.com/complynx/rpssl4bu/backend/pkg"
	"github.com/complynx/rpssl4bu/backend/pkg/types"
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
	{types.Rock, types.Rock}:     false,
	{types.Paper, types.Rock}:    true,
	{types.Scissors, types.Rock}: false,
	{types.Lizard, types.Rock}:   false,
	{types.Spock, types.Rock}:    true,

	// against paper
	{types.Rock, types.Paper}:     false,
	{types.Paper, types.Paper}:    false,
	{types.Scissors, types.Paper}: true,
	{types.Lizard, types.Paper}:   true,
	{types.Spock, types.Paper}:    false,

	// against scissors
	{types.Rock, types.Scissors}:     true,
	{types.Paper, types.Scissors}:    false,
	{types.Scissors, types.Scissors}: false,
	{types.Lizard, types.Scissors}:   false,
	{types.Spock, types.Scissors}:    true,

	// against lizard
	{types.Rock, types.Lizard}:     true,
	{types.Paper, types.Lizard}:    false,
	{types.Scissors, types.Lizard}: true,
	{types.Lizard, types.Lizard}:   false,
	{types.Spock, types.Lizard}:    false,

	// against spock
	{types.Rock, types.Spock}:     false,
	{types.Paper, types.Spock}:    true,
	{types.Scissors, types.Spock}: false,
	{types.Lizard, types.Spock}:   true,
	{types.Spock, types.Spock}:    false,
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
	for i := types.Rock; i <= types.Spock; i++ {
		ret = append(ret, i)
	}
	return ret, nil
}

func (g *game) Choice(ctx context.Context) (types.Choice, error) {
	num, err := g.rng.Rand(ctx)
	if err != nil {
		return types.Lizard, fmt.Errorf("generate random number: %w", err)
	}

	return types.IntToChoiceErr((num % 5) + 1)
}

func (g *game) Play(ctx context.Context, player types.Choice) (types.Result, types.Choice, error) {

	computerChoice, err := g.Choice(ctx)
	if err != nil {
		return types.Tie, types.Spock, fmt.Errorf("get computer choice: %w", err)
	}

	return GameResult(player, computerChoice), computerChoice, nil
}
