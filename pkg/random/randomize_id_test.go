package random

import (
	"context"
	"errors"
	"testing"

	"github.com/complynx/rpssl4bu/pkg/mocks"
	"github.com/complynx/rpssl4bu/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRandomizeID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		rng := mocks.NewRandomProvider(t)
		defer rng.AssertExpectations(t)

		rng.EXPECT().Rand(mock.Anything).Times(12).Return(0, nil)

		gi, err := RandomID(context.Background(), rng)
		assert.NoError(t, err)
		assert.Equal(t, types.GameID(0xd4a51000d4a51000), gi)
	})
	t.Run("ok 99", func(t *testing.T) {
		rng := mocks.NewRandomProvider(t)
		defer rng.AssertExpectations(t)

		rng.EXPECT().Rand(mock.Anything).Times(12).Return(99, nil)

		gi, err := RandomID(context.Background(), rng)
		assert.NoError(t, err)
		assert.Equal(t, types.GameID(0xa94a1fffa94a1fff), gi)
	})
	t.Run("ok 55", func(t *testing.T) {
		rng := mocks.NewRandomProvider(t)
		defer rng.AssertExpectations(t)

		rng.EXPECT().Rand(mock.Anything).Times(12).Return(55, nil)

		gi, err := RandomID(context.Background(), rng)
		assert.NoError(t, err)
		assert.Equal(t, types.GameID(0x2e5618e32e5618e3), gi)
	})
	t.Run("fail in 1st cycle", func(t *testing.T) {
		rng := mocks.NewRandomProvider(t)
		defer rng.AssertExpectations(t)

		rng.EXPECT().Rand(mock.Anything).Times(1).Return(55, nil)
		rng.EXPECT().Rand(mock.Anything).Times(1).Return(0, errors.New("test"))

		_, err := RandomID(context.Background(), rng)
		assert.EqualError(t, err, "get random number from rng: test")
	})
	t.Run("fail in 2nd cycle", func(t *testing.T) {
		rng := mocks.NewRandomProvider(t)
		defer rng.AssertExpectations(t)

		rng.EXPECT().Rand(mock.Anything).Times(8).Return(55, nil)
		rng.EXPECT().Rand(mock.Anything).Times(1).Return(0, errors.New("test"))

		_, err := RandomID(context.Background(), rng)
		assert.EqualError(t, err, "get random number from rng: test")
	})
}
