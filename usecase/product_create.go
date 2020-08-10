package usecase

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type ProductCreateRequest struct {
	Name  string
	Price int
}

type ProductCreateResponse struct {
	ID    int
	Name  string
	Price int
}

type ProductCreate struct {
	repo struct {
		product repository.Product
	}
}

func NewProductCreate(productRepo repository.Product) *ProductCreate {
	return &ProductCreate{
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductCreate) Create(ctx context.Context, req ProductCreateRequest) (*ProductCreateResponse, error) {
	id, err := p.repo.product.NextID(ctx)
	if err != nil {
		return nil, err
	}

	product := model.NewProduct(id, req.Name, req.Price)
	if err := p.repo.product.Store(ctx, *product); err != nil {
		return nil, err
	}

	return &ProductCreateResponse{
		ID:    int(product.ID()),
		Name:  req.Name,
		Price: req.Price,
	}, nil
}
