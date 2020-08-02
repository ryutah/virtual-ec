package usecase

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type ProductAddRequest struct {
	Name  string
	Price int
}

type ProductAddResponse struct {
	ID    int
	Name  string
	Price int
}

type ProductCreator struct {
	repo struct {
		product repository.Product
	}
}

func NewProductCreator(productRepo repository.Product) *ProductCreator {
	return &ProductCreator{
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductCreator) Append(ctx context.Context, req ProductAddRequest) (*ProductAddResponse, error) {
	id, err := p.repo.product.NextID(ctx)
	if err != nil {
		return nil, err
	}

	product := model.NewProduct(id, req.Name, req.Price)
	if err := p.repo.product.Store(ctx, *product); err != nil {
		return nil, err
	}

	return &ProductAddResponse{
		ID:    int(product.ID()),
		Name:  req.Name,
		Price: req.Price,
	}, nil
}
