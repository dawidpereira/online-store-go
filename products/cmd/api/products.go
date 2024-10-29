package main

import (
	"log"
	"net/http"
	"products/internal/store"
)

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var createProductRequest CreateProductRequest

	if err := readJSON(w, r, &createProductRequest); err != nil {
		err := writeJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	product := &store.Product{
		Name:        createProductRequest.Name,
		Description: createProductRequest.Description,
		Category:    createProductRequest.Category,
	}

	if err := app.store.Products.Create(product); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}
