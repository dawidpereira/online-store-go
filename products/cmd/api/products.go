package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"products/internal/store"
	"strconv"
)

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var createProductRequest CreateProductRequest

	if err := readJSON(w, r, &createProductRequest); err != nil {
		app.badRequestError(w, r, err)

		return
	}

	product := &store.Product{
		Name:        createProductRequest.Name,
		Description: createProductRequest.Description,
		Category:    createProductRequest.Category,
	}

	if err := app.store.Products.Create(product); err != nil {
		app.internalServerError(w, r, err)

		return
	}

	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	var updateProductRequest CreateProductRequest

	if err := readJSON(w, r, &updateProductRequest); err != nil {
		app.badRequestError(w, r, err)

		return
	}

	productForm := &store.Product{
		Name:        updateProductRequest.Name,
		Description: updateProductRequest.Description,
		Category:    updateProductRequest.Category,
	}

	product, err := app.store.Products.Update(id, productForm)
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.store.Products.List()
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	if err := writeJSON(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	product, err := app.store.Products.Get(id)
	if err != nil {
		app.notFoundError(w, r)

		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	if err := app.store.Products.Delete(id); err != nil {
		app.internalServerError(w, r, err)

		return
	}

	if err := writeJSON(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
