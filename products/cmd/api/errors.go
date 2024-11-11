package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "path", r.URL.Path, "error", err.Error())
	err = writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
	if err != nil {
		app.logger.Fatal(err)
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error", "path", r.URL.Path, "error", err.Error())
	err = writeJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		app.logger.Fatal(err)
	}
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	app.logger.Errorw("not found error", "path", r.URL.Path)
	err := writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
	if err != nil {
		app.logger.Fatal(err)
	}
}
