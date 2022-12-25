package storage

import (
	"testing"

	"github.com/complynx/rpssl4bu/backend/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestSimpleStorage(t *testing.T) {
	capacity := 3
	s := NewSimple(capacity)

	// Test GetLastScores on empty storage
	scores, err := s.GetLastScores()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(scores))

	// Test SetLastScore and GetLastScores
	err = s.SetLastScore(types.Win)
	assert.NoError(t, err)
	scores, err = s.GetLastScores()
	assert.NoError(t, err)
	assert.Equal(t, []types.Result{types.Win}, scores)

	// Test SetLastScore and GetLastScores with capacity
	err = s.SetLastScore(types.Lose)
	assert.NoError(t, err)
	err = s.SetLastScore(types.Tie)
	assert.NoError(t, err)
	scores, err = s.GetLastScores()
	assert.NoError(t, err)
	assert.Equal(t, []types.Result{types.Tie, types.Lose, types.Win}, scores)

	// Test SetLastScore with overflow
	err = s.SetLastScore(types.Tie)
	assert.NoError(t, err)
	scores, err = s.GetLastScores()
	assert.NoError(t, err)
	assert.Equal(t, []types.Result{types.Tie, types.Tie, types.Lose}, scores)

	// Test ClearScores
	err = s.ClearScores()
	assert.NoError(t, err)
	scores, err = s.GetLastScores()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(scores))
}
