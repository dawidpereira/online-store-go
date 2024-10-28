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
	products map[int64]*Product
	nextID   int64
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[int64]*Product),
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

	s.products[product.ID] = product
	return nil
}

func (s *ProductStore) Get(id int64) (*Product, error) {
	s.Lock()
	defer s.Unlock()

	product, exists := s.products[id]
	if !exists {
		return nil, fmt.Errorf("product with id %v not found", id)
	}
	return product, nil
}

func (s *ProductStore) Update(id int64, updatedProduct *Product) error {
	s.Lock()
	defer s.Unlock()

	product, exists := s.products[id]
	if !exists {
		return fmt.Errorf("product with id %v not found", id)
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Category = updatedProduct.Category
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	s.products[id] = product
	return nil
}

func (s *ProductStore) Delete(id int64) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.products[id]; !exists {
		return fmt.Errorf("product with id %v not found", id)
	}
	delete(s.products, id)
	return nil
}
