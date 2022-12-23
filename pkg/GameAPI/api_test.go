package gameapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/complynx/rpssl4bu/pkg/mocks"
	"github.com/complynx/rpssl4bu/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHTTPCode(t *testing.T) {
	// Test cases
	testCases := []struct {
		code int
		text string
	}{
		{200, "OK"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
	}

	for _, tc := range testCases {
		// Setup
		w := httptest.NewRecorder()

		// Test
		httpCode(w, tc.code)

		// Assert
		assert.Equal(t, tc.code, w.Code, "Wrong status code")
		assert.Equal(t, tc.text, http.StatusText(tc.code), "Wrong status text")
	}
}

func TestChoices(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		game           *mocks.Game
		request        *http.Request
		requestChoices []types.Choice
		expectedStatus int
		expectedBody   interface{}
		expectedLogs   []string
		expectedErr    error
	}{
		{
			name:           "Choices success",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodGet, "/choices", nil),
			requestChoices: []types.Choice{types.Lizard, types.Paper},
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`[{"id":3,"name":"lizard"},{"id":1,"name":"paper"}]`),
			expectedLogs:   []string{},
		},
		{
			name:           "Choices error",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodGet, "/choices", nil),
			expectedStatus: http.StatusInternalServerError,
			expectedLogs:   []string{"Error during request processing"},
			expectedErr:    errors.New("Error getting choices"),
		},
		// ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)
			// Create gameAPI instance
			api := NewGameAPI(tc.game, observedLogger)

			// Test
			w := httptest.NewRecorder()
			tc.game.On("Choices", mock.Anything).Return(tc.requestChoices, tc.expectedErr)
			api.Choices(w, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong status code")
			if tc.expectedBody != nil {
				assert.Equal(t, tc.expectedBody, w.Body.Bytes(), "Wrong response body")
			}
			// Assert logs
			logs := observedLogs.All()
			for _, log := range tc.expectedLogs {
				found := false
				for _, entry := range logs {
					if strings.Contains(entry.Message, log) {
						found = true
						break
					}
				}
				assert.True(t, found, "Log message not found: "+log)
			}

			// Assert mock expectations
			tc.game.AssertExpectations(t)
		})
	}
}

func TestChoice(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		game           *mocks.Game
		request        *http.Request
		expectedStatus int
		expectedBody   []byte
		expectedLogs   []string
		expectedErr    error
	}{
		{
			name:           "Choice success",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodGet, "/choice", nil),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"id":0,"name":"rock"}`),
			expectedLogs:   []string{},
			expectedErr:    nil,
		},
		{
			name:           "Choice error",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodGet, "/choice", nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			expectedLogs:   []string{"Error during request processing"},
			expectedErr:    errors.New("Error getting choice"),
		},
		// ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)
			// Create gameAPI instance
			api := NewGameAPI(tc.game, observedLogger)

			// Test
			w := httptest.NewRecorder()
			tc.game.On("Choice", mock.Anything).Return(types.Rock, tc.expectedErr)
			api.Choice(w, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong status code")
			if tc.expectedBody != nil {
				assert.Equal(t, tc.expectedBody, w.Body.Bytes(), "Wrong response body")
			}
			// Assert logs
			logs := observedLogs.All()
			for _, log := range tc.expectedLogs {
				found := false
				for _, entry := range logs {
					if strings.Contains(entry.Message, log) {
						found = true
						break
					}
				}
				assert.True(t, found, "Log message not found: "+log)
			}

			// Assert mock expectations
			tc.game.AssertExpectations(t)
		})
	}
}

func TestPlay(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		game           *mocks.Game
		request        *http.Request
		expectedStatus int
		expectedBody   []byte
		expectedLogs   []string
		expectedErr    error
	}{
		{
			name:           "Play success",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodPost, "/play", strings.NewReader("{\"player\":3}")),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"results":"win","player":3,"computer":3}`),
			expectedLogs:   []string{},
		},
		{
			name:           "Play error",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodPost, "/play", strings.NewReader("{\"player\":3}")),
			expectedStatus: http.StatusInternalServerError,
			expectedLogs:   []string{"Error during request processing"},
			expectedErr:    errors.New("Error getting play"),
		},
		// ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)
			// Create gameAPI instance
			api := NewGameAPI(tc.game, observedLogger)

			// Test
			w := httptest.NewRecorder()
			if tc.expectedBody != nil {
				tc.game.On("Play", mock.Anything, types.Lizard).Return(types.Win, types.Lizard, tc.expectedErr)
			} else {
				tc.game.On("Play", mock.Anything, types.Lizard).Return(types.Tie, types.Lizard, tc.expectedErr)

			}
			api.Play(w, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong status code")
			if tc.expectedBody != nil {
				assert.Equal(t, tc.expectedBody, w.Body.Bytes(), "Wrong response body")
			}
			// Assert logs
			logs := observedLogs.All()
			for _, log := range tc.expectedLogs {
				found := false
				for _, entry := range logs {
					if strings.Contains(entry.Message, log) {
						found = true
						break
					}
				}
				assert.True(t, found, "Log message not found: "+log)
			}

			// Assert mock expectations
			tc.game.AssertExpectations(t)
		})
	}
}

func TestMethods(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		game           *mocks.Game
		request        *http.Request
		expectedStatus int
	}{
		{
			name:           "Play",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodGet, "/play", nil),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Choice",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodPost, "/choice", nil),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Choices",
			game:           mocks.NewGame(t),
			request:        httptest.NewRequest(http.MethodPost, "/choices", nil),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		// ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)
			// Create gameAPI instance
			api := NewGameAPI(tc.game, observedLogger)

			// Test
			w := httptest.NewRecorder()
			switch tc.request.URL.Path {
			case "/choices":
				api.Choices(w, tc.request)
			case "/choice":
				api.Choice(w, tc.request)
			case "/play":
				api.Play(w, tc.request)
			}

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong status code")
			// Assert logs
			assert.Equal(t, 0, observedLogs.Len())

			// Assert mock expectations
			tc.game.AssertExpectations(t)
		})
	}
}

func TestGameAPI_Play_BadRequest(t *testing.T) {
	// Setup
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	game := mocks.NewGame(t)
	api := NewGameAPI(game, observedLogger)

	// Test
	request, err := http.NewRequest(http.MethodPost, "/play", strings.NewReader("{invalid json}"))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	api.Play(w, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code, "Wrong status code")

	// Assert logs
	assert.Equal(t, 0, observedLogs.Len(), "Wrong number of logs")
}
