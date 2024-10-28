package store

type Storage struct {
	Products interface {
		Create(product *Product) error
		Get(id int64) (*Product, error)
		Update(id int64, updatedProduct *Product) error
		Delete(id int64) error
	}
}

func NewStorage() Storage {
	return Storage{
		Products: NewProductStore(),
	}
}