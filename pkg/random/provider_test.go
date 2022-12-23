package random

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProvider_Rand(t *testing.T) {
	// Test cases
	testCases := []struct {
		name         string
		provider     *provider
		response     *response
		expectedRand int
		expectedErr  error
	}{
		{
			name: "Success",
			provider: &provider{
				addr: "http://example.com",
			},
			response: &response{
				RandomNumber: 43,
			},
			expectedRand: 42,
			expectedErr:  nil,
		},
		{
			name: "Bad request",
			provider: &provider{
				addr: "http://example.com",
			},
			response:     nil,
			expectedRand: 0,
			expectedErr:  fmt.Errorf("unmarshal body: EOF"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.response != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(tc.response)
				} else {
					http.Error(w, "", http.StatusBadRequest)
				}
			}))
			defer server.Close()

			tc.provider.addr = server.URL

			rand, err := tc.provider.Rand(context.Background())
			assert.Equal(t, tc.expectedRand, rand, "Wrong random number")
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error(), "Wrong error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
		})
	}
}

func TestProvider_Rand_timeout(t *testing.T) {
	// Setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(RequestTimeout + time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response{
			RandomNumber: 42,
		})
	}))
	defer server.Close()
	p := &provider{
		addr: server.URL,
	}

	rand, err := p.Rand(context.Background())
	assert.Equal(t, 0, rand, "Wrong random number")
	assert.Error(t, err, "Error expected")
	assert.Regexp(t, regexp.MustCompile("send request: Get \"http://127.0.0.1:[0-9]+\": context deadline exceeded"), err.Error(), "Wrong error")
}

func TestProvider_Rand_badRequest(t *testing.T) {
	// Setup

	p := &provider{
		addr: "\b",
	}

	rand, err := p.Rand(context.Background())
	assert.Equal(t, 0, rand, "Wrong random number")
	assert.EqualError(t, err, "create request: parse \"\\b\": net/url: invalid control character in URL")
}
