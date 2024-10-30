package store

import (
	"fmt"
	"sync"
	"time"
)

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ProductStore TODO: Change implementation to use a database
type ProductStore struct {
	sync.Mutex
	products []*Product
	nextID   int64
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make([]*Product, 0),
		nextID:   1,
	}
}

func (s *ProductStore) Create(product *Product) error {
	s.Lock()
	defer s.Unlock()

	product.ID = s.nextID
	s.nextID++

	currentTime := time.Now().Format(time.RFC3339)
	product.CreatedAt = currentTime
	product.UpdatedAt = currentTime

	s.products = append(s.products, product)
	return nil
}

func (s *ProductStore) List() ([]*Product, error) {
	s.Lock()
	defer s.Unlock()

	return s.products, nil
}

func (s *ProductStore) Get(id int64) (*Product, error) {
	s.Lock()
	defer s.Unlock()

	product, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})

	if !exists {
		return nil, fmt.Errorf("product with id %v not found", id)
	}

	return product, nil
}

func (s *ProductStore) Update(id int64, updatedProduct *Product) (*Product, error) {
	s.Lock()
	defer s.Unlock()

	product, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})

	if !exists {
		return nil, fmt.Errorf("product with id %v not found", id)
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Category = updatedProduct.Category
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	return product, nil
}

func (s *ProductStore) Delete(id int64) error {
	s.Lock()
	defer s.Unlock()

	_, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})
	if !exists {
		return fmt.Errorf("product with id %v not found", id)
	}

	s.products = remove(s.products, func(product *Product) bool {
		return product.ID == id
	})

	return nil
}

func filter(products []Product, predicate func(product Product) bool) []Product {
	var filtered []Product
	for _, product := range products {
		if predicate(product) {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

func find(products []*Product, predicate func(product *Product) bool) (*Product, bool) {
	for _, product := range products {
		if predicate(product) {
			return product, true
		}
	}

	return nil, false
}

func remove(products []*Product, predicate func(product *Product) bool) []*Product {
	var filtered []*Product
	for _, product := range products {
		if !predicate(product) {
			filtered = append(filtered, product)
		}
	}

	return filtered
}
