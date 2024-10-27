package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := map[string]string{
		"status": "available",
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err) //TODO: Handle server error
	}
}
