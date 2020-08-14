package model

import (
	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
)

type NewProductInput struct {
	ID    ProductID `validate:"required"`
	Name  string    `validate:"required,len=255"`
	Price int       `validate:"required,gt=0"`
}

type ProductID int

type Product struct {
	id    ProductID
	name  string
	price int
}

func NewProduct(id ProductID, name string, price int) (*Product, error) {
	if name == "" {
		return nil, errors.Wrap(domain.ErrInvalidInput, "name should not be blank")
	}

	return &Product{
		id:    id,
		name:  name,
		price: price,
	}, nil
}

func ReCreateProduct(id ProductID, name string, price int) *Product {
	return &Product{
		id:    id,
		name:  name,
		price: price,
	}
}

func (p *Product) ID() ProductID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() int {
	return p.price
}

func (p *Product) NewReview(id ReviewID) *Review {
	return &Review{
		id:       id,
		reviewTo: p.ID(),
	}
}
