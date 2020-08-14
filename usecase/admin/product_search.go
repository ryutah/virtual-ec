package admin

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/ryutah/virtual-ec/lib/xlog"
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
	output ProductSearchOutputPort
	repo   struct {
		product repository.Product
	}
}

func NewProductSearch(output ProductSearchOutputPort, productRepo repository.Product) *ProductSearch {
	return &ProductSearch{
		output: output,
		repo: struct{ product repository.Product }{
			product: productRepo,
		},
	}
}

func (p *ProductSearch) Search(ctx context.Context, input ProductSearchInputPort) (success bool) {
	result, err := p.repo.product.Search(ctx, p.toQuery(input))
	if err != nil {
		return p.handleError(ctx, err)
	}

	var products []*ProductSearchProductDetail
	for _, product := range result.Products {
		products = append(products, &ProductSearchProductDetail{
			ID:    int(product.ID()),
			Name:  product.Name(),
			Price: product.Price(),
		})
	}
	p.output.Success(ProductSearchSuccess{
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

func (p *ProductSearch) handleError(ctx context.Context, err error) bool {
	xlog.Errorf(ctx, "failed to product search: %+v", err)
	p.output.Failed(ProductSearchFailed{
		Err: productSearchErrorMessages.failed(),
	})
	return false
}
