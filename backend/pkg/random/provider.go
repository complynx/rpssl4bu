package random

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/complynx/rpssl4bu/backend/pkg"
)

const RequestTimeout = 1 * time.Second

type provider struct {
	addr string
}

type response struct {
	RandomNumber int `json:"random_number"`
}

func NewProvider(addr string) pkg.RandomProvider {
	return &provider{
		addr: addr,
	}
}

func (p *provider) Rand(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.addr, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("unmarshal body: %w", err)
	}

	return res.RandomNumber - 1, nil
}
