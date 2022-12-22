package game

import (
	"context"
	"errors"
	"testing"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/mocks"
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
		tie, win := GameResult(0, 0)
		s.True(tie)
		s.False(win)
	})
	s.Run("win", func() {
		tie, win := GameResult(1, 0)
		s.True(win)
		s.False(tie)
	})
	s.Run("lose", func() {
		tie, win := GameResult(0, 1)
		s.False(win)
		s.False(tie)
	})
}

func (s *gameTestSuite) TestStringResult() {
	s.Run("tie", func() {
		res := StringResult(0, 0)
		s.Equal("tie", res)
	})
	s.Run("win", func() {
		res := StringResult(1, 0)
		s.Equal("win", res)
	})
	s.Run("lose", func() {
		res := StringResult(0, 1)
		s.Equal("lose", res)
	})
}

func (s *gameTestSuite) TestChoices() {
	game := NewGame(nil)

	res, err := game.Choices(context.Background())
	s.NoError(err)
	s.Equal([]pkg.Choice{
		{ID: 0, Name: "rock"},
		{ID: 1, Name: "paper"},
		{ID: 2, Name: "scissors"},
		{ID: 3, Name: "lizard"},
		{ID: 4, Name: "spock"},
	}, res)
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
		s.Equal(pkg.Choice{
			ID:   2,
			Name: "scissors",
		}, res)
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
		s.Equal(pkg.Choice{
			ID:   2,
			Name: "scissors",
		}, choice)
		s.Equal(res, "lose")
	})
}
