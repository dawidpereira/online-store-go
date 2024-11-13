package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetProduct(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	t.Run("should return a product", func(t *testing.T) {
		// Arrange
		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusOK, rr.Code)

	})

	t.Run("should return not found when there is no product", func(t *testing.T) {
		// Arrange
		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusNotFound, rr.Code)
	})
}

func TestListProduct(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	t.Run("should return a product", func(t *testing.T) {
		// Arrange
		req, err := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		response := decodeResponseBody(t, rr.Result())
		data, ok := response.Data.([]interface{})
		if !ok {
			t.Fatalf("expected a slice of interface{}, got %T", response.Data)
		}
		expectedCount := 10
		if len(data) != expectedCount {
			t.Errorf("expected %d product, got %d", expectedCount, len(data))
		}
	})
	//TODO: Test filters, pagination, and ordering
}

func TestCreateProduct(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	t.Run("should create a product", func(t *testing.T) {
		// Arrange
		product := CreateProductRequest{
			Name:        "Test product",
			Category:    "Category",
			Description: "Description",
		}

		body, err := json.Marshal(product)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusCreated, rr.Code)
	})
}

func TestUpdateProduct(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	//TODO: Verify if properties are mapped correctly
	t.Run("should update a product", func(t *testing.T) {
		// Arrange
		product := UpdateProductRequest{
			Name:        "Test product",
			Category:    "Category",
			Description: "Description",
		}

		body, err := json.Marshal(product)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should return not found when there is no product", func(t *testing.T) {
		// Arrange
		product := UpdateProductRequest{
			Name:        "Test product",
			Category:    "Category",
			Description: "Description",
		}

		body, err := json.Marshal(product)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, "/api/v1/products/999", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusNotFound, rr.Code)
	})
}

func TestDeleteProduct(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	t.Run("should delete a product", func(t *testing.T) {
		// Arrange
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		rr := executeRequest(req, mux)

		// Assert
		assertResponseCode(t, http.StatusNoContent, rr.Code)
	})
}
