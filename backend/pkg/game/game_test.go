package game

import (
	"context"
	"errors"
	"testing"

	"github.com/complynx/rpssl4bu/backend/pkg/mocks"
	"github.com/complynx/rpssl4bu/backend/pkg/types"
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
		res := GameResult(types.Paper, types.Paper)
		s.Equal(types.Tie, res)
	})
	s.Run("win", func() {
		res := GameResult(types.Paper, types.Rock)
		s.Equal(types.Win, res)
	})
	s.Run("lose", func() {
		res := GameResult(types.Paper, types.Scissors)
		s.Equal(types.Lose, res)
	})
	s.Run("VS", func() {
		s.Run("Rock vs Rock", func() {
			res := GameResult(types.Rock, types.Rock)
			s.Equal(types.Tie, res)
		})
		s.Run("Rock vs Paper", func() {
			res := GameResult(types.Rock, types.Paper)
			s.Equal(types.Lose, res)
		})
		s.Run("Rock vs Scissors", func() {
			res := GameResult(types.Rock, types.Scissors)
			s.Equal(types.Win, res)
		})
		s.Run("Rock vs Lizard", func() {
			res := GameResult(types.Rock, types.Lizard)
			s.Equal(types.Win, res)
		})
		s.Run("Rock vs Spock", func() {
			res := GameResult(types.Rock, types.Spock)
			s.Equal(types.Lose, res)
		})

		s.Run("Paper vs Rock", func() {
			res := GameResult(types.Paper, types.Rock)
			s.Equal(types.Win, res)
		})
		s.Run("Paper vs Paper", func() {
			res := GameResult(types.Paper, types.Paper)
			s.Equal(types.Tie, res)
		})
		s.Run("Paper vs Scissors", func() {
			res := GameResult(types.Paper, types.Scissors)
			s.Equal(types.Lose, res)
		})
		s.Run("Paper vs Lizard", func() {
			res := GameResult(types.Paper, types.Lizard)
			s.Equal(types.Lose, res)
		})
		s.Run("Paper vs Spock", func() {
			res := GameResult(types.Paper, types.Spock)
			s.Equal(types.Win, res)
		})

		s.Run("Scissors vs Rock", func() {
			res := GameResult(types.Scissors, types.Rock)
			s.Equal(types.Lose, res)
		})
		s.Run("Scissors vs Paper", func() {
			res := GameResult(types.Scissors, types.Paper)
			s.Equal(types.Win, res)
		})
		s.Run("Scissors vs Scissors", func() {
			res := GameResult(types.Scissors, types.Scissors)
			s.Equal(types.Tie, res)
		})
		s.Run("Scissors vs Lizard", func() {
			res := GameResult(types.Scissors, types.Lizard)
			s.Equal(types.Win, res)
		})
		s.Run("Scissors vs Spock", func() {
			res := GameResult(types.Scissors, types.Spock)
			s.Equal(types.Lose, res)
		})

		s.Run("Lizard vs Rock", func() {
			res := GameResult(types.Lizard, types.Rock)
			s.Equal(types.Lose, res)
		})
		s.Run("Lizard vs Paper", func() {
			res := GameResult(types.Lizard, types.Paper)
			s.Equal(types.Win, res)
		})
		s.Run("Lizard vs Scissors", func() {
			res := GameResult(types.Lizard, types.Scissors)
			s.Equal(types.Lose, res)
		})
		s.Run("Lizard vs Lizard", func() {
			res := GameResult(types.Lizard, types.Lizard)
			s.Equal(types.Tie, res)
		})
		s.Run("Lizard vs Spock", func() {
			res := GameResult(types.Lizard, types.Spock)
			s.Equal(types.Win, res)
		})

		s.Run("Spock vs Rock", func() {
			res := GameResult(types.Spock, types.Rock)
			s.Equal(types.Win, res)
		})
		s.Run("Spock vs Paper", func() {
			res := GameResult(types.Spock, types.Paper)
			s.Equal(types.Lose, res)
		})
		s.Run("Spock vs Scissors", func() {
			res := GameResult(types.Spock, types.Scissors)
			s.Equal(types.Win, res)
		})
		s.Run("Spock vs Lizard", func() {
			res := GameResult(types.Spock, types.Lizard)
			s.Equal(types.Lose, res)
		})
		s.Run("Spock vs Spock", func() {
			res := GameResult(types.Spock, types.Spock)
			s.Equal(types.Tie, res)
		})
	})
}

func (s *gameTestSuite) TestChoices() {
	game := NewGame(nil)

	res, err := game.Choices(context.Background())
	s.NoError(err)
	s.Equal([]types.Choice{1, 2, 3, 4, 5}, res)
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

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(1, errors.New("test"))

		game := NewGame(rng)

		_, _, err := game.Play(context.Background(), 1)
		s.EqualError(err, "get computer choice: generate random number: test")
	})
	s.Run("ok", func() {
		rng := mocks.NewRandomProvider(s.T())
		defer rng.AssertExpectations(s.T())

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(2, nil)

		game := NewGame(rng)

		res, choice, err := game.Play(context.Background(), 2)
		s.NoError(err)
		s.Equal(types.Scissors, choice)
		s.Equal(types.Lose, res)
	})
}
