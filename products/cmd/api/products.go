package main

import (
	"github.com/go-chi/chi/v5"
	"log"
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

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		err := writeJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}

	var updateProductRequest CreateProductRequest

	if err := readJSON(w, r, &updateProductRequest); err != nil {
		err := writeJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	productForm := &store.Product{
		Name:        updateProductRequest.Name,
		Description: updateProductRequest.Description,
		Category:    updateProductRequest.Category,
	}

	product, err := app.store.Products.Update(id, productForm)
	if err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.store.Products.List()
	if err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := writeJSON(w, http.StatusOK, products); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		err := writeJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}

	product, err := app.store.Products.Get(id)
	if err != nil {
		err := writeJSONError(w, http.StatusNotFound, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		err := writeJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := app.store.Products.Delete(id); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := writeJSON(w, http.StatusNoContent, nil); err != nil {
		err := writeJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}
