package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("interlan server error: %s path: %s", err, r.URL.Path)
	err = writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s", err, r.URL.Path)
	err = writeJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	log.Printf("not found error: path: %s", r.URL.Path)
	err := writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
	if err != nil {
		log.Fatal(err)
	}
}
