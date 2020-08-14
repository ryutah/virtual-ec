package admin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/ryutah/virtual-ec/lib/xlog"
)

var productFindFailedErrorMessages = struct {
	notFound func(model.ProductID) string
	failed   func(model.ProductID) string
}{
	notFound: func(id model.ProductID) string { return fmt.Sprintf("Product(%v)は存在しません", id) },
	failed:   func(id model.ProductID) string { return fmt.Sprintf("Product(%v)の取得に失敗しました", id) },
}

type ProductFindOutputPort interface {
	Success(ProductFindSuccess)
	NotFound(ProductFindFailed)
	Failed(ProductFindFailed)
}

type (
	ProductFindFailed struct {
		Err string
	}

	ProductFindSuccess struct {
		ID    int
		Name  string
		Price int
	}
)

type ProductFind struct {
	output ProductFindOutputPort
	repo   struct {
		product repository.Product
	}
}

func NewProductFind(output ProductFindOutputPort, productRepo repository.Product) *ProductFind {
	return &ProductFind{
		output: output,
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductFind) Find(ctx context.Context, id int) (success bool) {
	product, err := p.repo.product.Get(ctx, model.ProductID(id))
	if err != nil {
		return p.handleError(ctx, model.ProductID(id), err)
	}
	p.output.Success(ProductFindSuccess{
		ID:    int(product.ID()),
		Name:  product.Name(),
		Price: product.Price(),
	})
	return true
}

func (p *ProductFind) handleError(ctx context.Context, id model.ProductID, err error) bool {
	if errors.Is(err, domain.ErrNoSuchEntity) {
		xlog.Warningf(ctx, "product not found: %+v", err)
		p.output.NotFound(ProductFindFailed{
			Err: productFindFailedErrorMessages.notFound(id),
		})
	} else {
		xlog.Errorf(ctx, "failed to find product: %+v", err)
		p.output.Failed(ProductFindFailed{
			Err: productFindFailedErrorMessages.failed(id),
		})
	}
	return false
}
