package admin

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/ryutah/virtual-ec/pkg/xlog"
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
		Err string
	}
)

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

func (p *ProductCreate) Create(ctx context.Context, in ProductCreateInputPort, out ProductCreateOutputPort) (success bool) {
	id, err := p.repo.product.NextID(ctx)
	if err != nil {
		return p.handleError(ctx, err, out)
	}

	product := model.ReCreateProduct(id, in.Name(), in.Price())
	if err := p.repo.product.Store(ctx, *product); err != nil {
		return p.handleError(ctx, err, out)
	}

	out.Success(ProductCreateSuccess{
		ID:    int(product.ID()),
		Name:  product.Name(),
		Price: product.Price(),
	})
	return true
}

func (p *ProductCreate) handleError(ctx context.Context, err error, out ProductCreateOutputPort) bool {
	xlog.Errorf(ctx, "failed to create product: %+v", err)
	out.Failed(ProductCreateFailed{
		Err: productCreateErrroMessages.failed(),
	})
	return false
}
