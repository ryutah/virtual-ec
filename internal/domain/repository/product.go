package repository

import (
	"context"

	"github.com/ryutah/virtual-ec/internal/domain/model"
)

type ProductQuery struct {
	name *string
}

func NewProductQuery() ProductQuery {
	return ProductQuery{}
}

func (p ProductQuery) WithName(name string) ProductQuery {
	p.name = &name
	return p
}

func (p ProductQuery) Name() (string, bool) {
	if p.name != nil {
		return *p.name, true
	}
	return "", false
}

type ProductSearchResult struct {
	Products []*model.Product
}

type Product interface {
	NextID(context.Context) (model.ProductID, error)
	Get(context.Context, model.ProductID) (*model.Product, error)
	Store(context.Context, model.Product) error
	Search(context.Context, ProductQuery) (*ProductSearchResult, error)
}
