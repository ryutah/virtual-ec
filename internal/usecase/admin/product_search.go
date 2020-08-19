package admin

import (
	"context"

	"github.com/ryutah/virtual-ec/internal/domain/repository"
	"github.com/ryutah/virtual-ec/pkg/xlog"
)

var productSearchErrorMessages = struct {
	failed func() string
}{
	failed: func() string { return "Productの検索に失敗しました" },
}

type (
	ProductSearchOutputPort interface {
		Success(ProductSearchSuccess)
		Failed(ProductSearchFailed)
	}

	ProductSearchInputPort interface {
		Name() string
	}
)

type (
	ProductSearchSuccess struct {
		Products []*ProductSearchProductDetail
	}

	ProductSearchProductDetail struct {
		ID    int
		Name  string
		Price int
	}

	ProductSearchFailed struct {
		Err string
	}
)

type ProductSearch struct {
	repo struct {
		product repository.Product
	}
}

func NewProductSearch(productRepo repository.Product) *ProductSearch {
	return &ProductSearch{
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductSearch) Search(ctx context.Context, in ProductSearchInputPort, out ProductSearchOutputPort) (success bool) {
	result, err := p.repo.product.Search(ctx, p.toQuery(in))
	if err != nil {
		return p.handleError(ctx, err, out)
	}

	var products []*ProductSearchProductDetail
	for _, product := range result.Products {
		products = append(products, &ProductSearchProductDetail{
			ID:    int(product.ID()),
			Name:  product.Name(),
			Price: product.Price(),
		})
	}
	out.Success(ProductSearchSuccess{
		Products: products,
	})
	return true
}

func (p *ProductSearch) toQuery(input ProductSearchInputPort) repository.ProductQuery {
	query := repository.NewProductQuery()
	if input.Name() != "" {
		query = query.WithName(input.Name())
	}
	return query
}

func (p *ProductSearch) handleError(ctx context.Context, err error, out ProductSearchOutputPort) bool {
	xlog.Errorf(ctx, "failed to product search: %+v", err)
	out.Failed(ProductSearchFailed{
		Err: productSearchErrorMessages.failed(),
	})
	return false
}
