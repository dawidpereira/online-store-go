package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status != http.StatusNoContent {
		return json.NewEncoder(w).Encode(data)
	}

	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	data := map[string]string{"error": message}

	return writeJSON(w, status, data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any, logger *zap.SugaredLogger) error {
	if r.Body == nil {
		logger.Fatal("request body is nil")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(r.Body)

	maxBytes := int64(1 << 20) // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
