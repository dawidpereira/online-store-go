package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"products/internal/store"
	"strconv"
)

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=100"`
	Category    string `json:"category" validate:"required,max=50"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var createProductRequest CreateProductRequest
	if err := readJSON(w, r, &createProductRequest); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(createProductRequest); err != nil {
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

type UpdateProductRequest struct {
	Name        string `json:"name" validator:"required,max=100"`
	Description string `json:"description" validator:"required,max=100"`
	Category    string `json:"category" validator:"required,max=50"`
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	var updateProductRequest UpdateProductRequest

	if err := readJSON(w, r, &updateProductRequest); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(updateProductRequest); err != nil {
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
	pq, err := store.ParseListProductsQuery(r)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	if err := Validate.Struct(pq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	products, err := app.store.Products.List(pq)
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	products.Next = pq.GetNextURL(r)
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
