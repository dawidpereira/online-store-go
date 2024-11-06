package store

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

type PaginatedQuery struct {
	Limit int   `json:"limit" validate:"required,gte=1,lte=50" default:"10"`
	Page  int   `json:"page" validate:"gte=0" default:"1"`
	Order Order `json:"order,omitempty" validate:"oneof=ASC DESC" default:"ASC"`
}

type PaginatedResponse struct {
	Limit int         `json:"limit" validate:"required,gte=1,lte=50" default:"10"`
	Page  int         `json:"page" validate:"gte=0" default:"1"`
	Order Order       `json:"order,omitempty" validate:"oneof=ASC DESC" default:"ASC"`
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
	Next  string      `json:"next,omitempty"`
}

type ListProductsQuery struct {
	PaginatedQuery `json:",inline"`
	Search         string   `json:"search"`
	Category       []string `json:"category"`
}

func (q *ListProductsQuery) GetNextURL(r *http.Request) string {
	nextURL := r.URL.Query()
	nextURL.Set("page", strconv.Itoa(q.Page+1))
	return nextURL.Encode()
}

func ParseListProductsQuery(r *http.Request) (ListProductsQuery, error) {
	query := ListProductsQuery{}

	paginatedQuery, err := ParsePaginatedQuery(r)
	if err != nil {
		return query, err
	}
	query.PaginatedQuery = paginatedQuery
	query.Search = r.URL.Query().Get("search")
	query.Category = r.URL.Query()["category"]

	return query, nil
}

func ParsePaginatedQuery(r *http.Request) (PaginatedQuery, error) {
	limitParam := chi.URLParam(r, "limit")
	if limitParam == "" {
		limitParam = "10"
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return PaginatedQuery{}, err
	}

	pageParam := chi.URLParam(r, "page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return PaginatedQuery{}, err
	}

	order := Order(r.URL.Query().Get("order"))
	if order == "" {
		order = ASC
	}

	query := PaginatedQuery{
		Limit: limit,
		Page:  page,
		Order: order,
	}

	return query, nil
}
