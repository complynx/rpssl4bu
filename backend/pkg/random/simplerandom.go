package random

import (
	"context"
	"math/rand"
	"time"

	"github.com/complynx/rpssl4bu/backend/pkg"
)

type simple struct{}

func NewSimpleRandom(addr string) pkg.RandomProvider {
	rand.Seed(time.Now().UnixNano())
	return &simple{}
}

func (p *simple) Rand(ctx context.Context) (int, error) {
	return rand.Intn(100), nil
}
