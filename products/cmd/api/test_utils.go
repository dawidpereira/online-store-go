package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"products/internal/store"
	"shared"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	storage := store.NewMockStorage()

	return &application{
		logger: logger,
		store:  storage,
		rateLimiter: shared.NewFixedWindowRateLimiter(shared.Config{
			RequestPerTimeFrame: 100,
			TimeFrame:           1,
			Enabled:             true,
		}, logger),
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func decodeResponseBody(t *testing.T, body *http.Response) store.PaginatedResponse {
	var response store.PaginatedResponse
	if err := json.NewDecoder(body.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	return response
}

func assertResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected status %d, got %d", expected, actual)
	}
}
