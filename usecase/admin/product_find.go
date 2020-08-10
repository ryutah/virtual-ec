package admin

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type ProductFindResponse struct {
	ID    int
	Name  string
	Price int
}

type ProductFind struct {
	repo struct {
		product repository.Product
	}
}

func NewProductFind(productRepo repository.Product) *ProductFind {
	return &ProductFind{
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductFind) Find(ctx context.Context, id int) (*ProductFindResponse, error) {
	product, err := p.repo.product.Get(ctx, model.ProductID(id))
	if err != nil {
		return nil, err
	}
	return &ProductFindResponse{
		ID:    int(product.ID()),
		Name:  product.Name(),
		Price: product.Price(),
	}, nil
}
