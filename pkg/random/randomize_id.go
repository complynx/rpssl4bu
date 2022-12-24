package random

import (
	"context"
	"fmt"

	"github.com/complynx/rpssl4bu/pkg"
	"github.com/complynx/rpssl4bu/pkg/types"
)

func RandomID(ctx context.Context, rng pkg.RandomProvider) (types.GameID, error) {
	var giHalf uint64
	var gi uint64

	for giHalf = 1; giHalf <= 0x0fffffffff; { // 0xf headroom for uniformity
		rn, err := rng.Rand(ctx)
		if err != nil {
			return types.GameID(gi), fmt.Errorf("get random number from rng: %w", err)
		}
		giHalf = 100*giHalf + uint64(rn)
	}
	gi = giHalf & 0xffffffff
	for giHalf = 1; giHalf <= 0x0fffffffff; { // 0xf headroom for uniformity
		rn, err := rng.Rand(ctx)
		if err != nil {
			return types.GameID(gi), fmt.Errorf("get random number from rng: %w", err)
		}
		giHalf = 100*giHalf + uint64(rn)
	}
	gi <<= 32
	gi |= giHalf & 0xffffffff
	return types.GameID(gi), nil
}
