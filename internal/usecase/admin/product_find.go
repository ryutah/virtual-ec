package admin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/internal/domain"
	"github.com/ryutah/virtual-ec/internal/domain/model"
	"github.com/ryutah/virtual-ec/internal/domain/repository"
	"github.com/ryutah/virtual-ec/pkg/xlog"
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

func (p *ProductFind) Find(ctx context.Context, id int, out ProductFindOutputPort) (success bool) {
	product, err := p.repo.product.Get(ctx, model.ProductID(id))
	if err != nil {
		return p.handleError(ctx, model.ProductID(id), err, out)
	}
	out.Success(ProductFindSuccess{
		ID:    int(product.ID()),
		Name:  product.Name(),
		Price: product.Price(),
	})
	return true
}

func (p *ProductFind) handleError(ctx context.Context, id model.ProductID, err error, out ProductFindOutputPort) bool {
	if errors.Is(err, domain.ErrNoSuchEntity) {
		xlog.Warningf(ctx, "product not found: %+v", err)
		out.NotFound(ProductFindFailed{
			Err: productFindFailedErrorMessages.notFound(id),
		})
	} else {
		xlog.Errorf(ctx, "failed to find product: %+v", err)
		out.Failed(ProductFindFailed{
			Err: productFindFailedErrorMessages.failed(id),
		})
	}
	return false
}
