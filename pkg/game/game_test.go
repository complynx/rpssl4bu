package game

import (
	"context"
	"errors"
	"testing"

	"github.com/complynx/rpssl4bu/pkg/mocks"
	"github.com/complynx/rpssl4bu/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestGame(t *testing.T) {
	suite.Run(t, new(gameTestSuite))
}

type gameTestSuite struct {
	suite.Suite
}

func (s *gameTestSuite) TestGameResult() {
	s.Run("tie", func() {
		res := GameResult(0, 0)
		s.Equal(types.Tie, res)
	})
	s.Run("win", func() {
		res := GameResult(1, 0)
		s.Equal(types.Win, res)
	})
	s.Run("lose", func() {
		res := GameResult(0, 1)
		s.Equal(types.Lose, res)
	})
}

func (s *gameTestSuite) TestChoices() {
	game := NewGame(nil)

	res, err := game.Choices(context.Background())
	s.NoError(err)
	s.Equal([]types.Choice{0, 1, 2, 3, 4}, res)
}

func (s *gameTestSuite) TestChoice() {
	s.Run("fail random", func() {
		rng := mocks.NewRandomProvider(s.T())
		defer rng.AssertExpectations(s.T())

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(0, errors.New("test"))

		game := NewGame(rng)

		_, err := game.Choice(context.Background())
		s.EqualError(err, "generate random number: test")
	})
	s.Run("ok", func() {
		rng := mocks.NewRandomProvider(s.T())
		defer rng.AssertExpectations(s.T())

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(2, nil)

		game := NewGame(rng)

		res, err := game.Choice(context.Background())
		s.NoError(err)
		s.Equal(types.Scissors, res)
	})
}
func (s *gameTestSuite) TestPlay() {
	s.Run("fail random", func() {
		rng := mocks.NewRandomProvider(s.T())
		defer rng.AssertExpectations(s.T())

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(0, errors.New("test"))

		game := NewGame(rng)

		_, _, err := game.Play(context.Background(), 1)
		s.EqualError(err, "get computer choice: generate random number: test")
	})
	s.Run("ok", func() {
		rng := mocks.NewRandomProvider(s.T())
		defer rng.AssertExpectations(s.T())

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(2, nil)

		game := NewGame(rng)

		res, choice, err := game.Play(context.Background(), 1)
		s.NoError(err)
		s.Equal(types.Scissors, choice)
		s.Equal(types.Lose, res)
	})
}
