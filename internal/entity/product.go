package entity

import (
	"errors"
	"time"

	"github.com/adrianomota/fullcycle/my-api/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("id od required")
	ErrInvalidId       = errors.New("invalid id")
	ErrNameIsRequire   = errors.New("id od required")
	ErrPriceIsRequired = errors.New("id od required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"Name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	err := product.Validate()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}

	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidId
	}

	if p.Name == "" {
		return ErrNameIsRequire
	}

	if p.Price == 0.0 {
		return ErrPriceIsRequired
	}

	if p.Price < 0.0 {
		return ErrInvalidPrice
	}

	return nil
}
