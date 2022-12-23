package storage

import (
	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/types"
)

type simple struct {
	scores   []types.Result
	capacity int
}

func NewSimple(capacity int) pkg.Storage {
	return &simple{
		capacity: capacity,
	}
}

// lists last scores
func (s *simple) GetLastScores() ([]types.Result, error) {
	l := len(s.scores)
	ret := make([]types.Result, l)
	for i := range s.scores {
		ret[l-i-1] = s.scores[i]
	}
	return ret, nil
}

// updates scoreboard adding last one, removing overflow if needed
func (s *simple) SetLastScore(r types.Result) error {
	s.scores = append(s.scores, r)
	if len(s.scores) > s.capacity {
		s.scores = s.scores[1:]
	}
	return nil
}

// cleares storage
func (s *simple) ClearScores() error {
	s.scores = nil
	return nil
}
