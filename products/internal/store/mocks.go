package store

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func NewMockStorage() Storage {
	return Storage{
		Products: NewMockProductStorage(),
	}
}

func NewMockProductStorage() *MockProductStore {
	store := &MockProductStore{
		products: make([]*Product, 0),
		nextID:   1,
	}

	for i := 1; i <= 10; i++ {
		_ = store.Create(&Product{
			Name:        fmt.Sprintf("Product %d", i),
			Description: fmt.Sprintf("Description for product %d", i),
			Category:    fmt.Sprintf("Category %d", i),
		})
	}

	return store
}

type MockProductStore struct {
	sync.Mutex
	products []*Product
	nextID   int64
}

func (s *MockProductStore) Create(product *Product) error {
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

func (s *MockProductStore) List(query ListProductsQuery) (PaginatedResponse, error) {
	s.Lock()
	defer s.Unlock()

	start := (query.Page - 1) * query.Limit
	end := start + query.Limit
	if end > len(s.products) {
		end = len(s.products)
	}

	if query.Order == DESC {
		start, end = len(s.products)-end, len(s.products)-start
	}

	filtered := filter(s.products, func(product *Product) bool {
		if query.Search != "" && !strings.Contains(product.Name, query.Search) {
			return false
		}
		if len(query.Category) > 0 && !contains(query.Category, product.Category) {
			return false
		}
		return true
	})

	return PaginatedResponse{
		Limit: query.Limit,
		Page:  query.Page,
		Order: query.Order,
		Total: len(s.products),
		Data:  filtered[start:end],
	}, nil
}

func (s *MockProductStore) Get(id int64) (*Product, error) {
	s.Lock()
	defer s.Unlock()

	product, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})

	if !exists {
		return nil, &ProductNotFoundError{ID: id}
	}

	return product, nil
}

func (s *MockProductStore) Update(id int64, updatedProduct *Product) (*Product, error) {
	s.Lock()
	defer s.Unlock()

	product, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})

	if !exists {
		return nil, &ProductNotFoundError{ID: id}
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Category = updatedProduct.Category
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	return product, nil
}

func (s *MockProductStore) Delete(id int64) error {
	s.Lock()
	defer s.Unlock()

	_, exists := find(s.products, func(product *Product) bool {
		return product.ID == id
	})

	if !exists {
		return &ProductNotFoundError{ID: id}
	}

	s.products = remove(s.products, func(product *Product) bool {
		return product.ID == id
	})

	return nil
}
