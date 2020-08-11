package admin

import (
	"context"
	"errors"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/ryutah/virtual-ec/lib/xlog"
)

var productCreateErrroMessages = struct {
	failed func() string
}{
	failed: func() string { return "Productの作成に失敗しました" },
}

type (
	ProductCreateOutputPort interface {
		Success(ProductCreateSuccess)
		Failed(ProductCreateFailed)
	}

	ProductCreateInputPort interface {
		Name() string
		Price() int
	}
)

type (
	ProductCreateSuccess struct {
		ID    int
		Name  string
		Price int
	}

	ProductCreateFailed struct {
		Err error
	}
)

type ProductCreate struct {
	output ProductCreateOutputPort
	repo   struct {
		product repository.Product
	}
}

func NewProductCreate(output ProductCreateOutputPort, productRepo repository.Product) *ProductCreate {
	return &ProductCreate{
		output: output,
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductCreate) Create(ctx context.Context, input ProductCreateInputPort) (success bool) {
	id, err := p.repo.product.NextID(ctx)
	if err != nil {
		return p.handleError(ctx, err)
	}

	product := model.NewProduct(id, input.Name(), input.Price())
	if err := p.repo.product.Store(ctx, *product); err != nil {
		return p.handleError(ctx, err)
	}

	p.output.Success(ProductCreateSuccess{
		ID:    int(product.ID()),
		Name:  product.Name(),
		Price: product.Price(),
	})
	return true
}

func (p *ProductCreate) handleError(ctx context.Context, err error) bool {
	xlog.Errorf(ctx, "failed to create product: %+v", err)
	p.output.Failed(ProductCreateFailed{
		Err: errors.New(productCreateErrroMessages.failed()),
	})
	return false
}
