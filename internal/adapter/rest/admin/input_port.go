package admin

import (
	"github.com/ryutah/virtual-ec/internal/adapter/rest/admin/internal"
	"github.com/ryutah/virtual-ec/internal/usecase/admin"
)

type productSearchInputPort struct {
	params internal.ProductSearchParams
}

var _ admin.ProductSearchInputPort = productSearchInputPort{}

func newProductSearchInputPort(p internal.ProductSearchParams) productSearchInputPort {
	return productSearchInputPort{
		params: p,
	}
}

func (p productSearchInputPort) Name() string {
	if p.params.Name == nil {
		return ""
	}
	return *p.params.Name
}

type productCreateInputPort struct {
	params internal.ProductCreateJSONRequestBody
}

var _ admin.ProductCreateInputPort = productCreateInputPort{}

func newProductCreateInputPort(p internal.ProductCreateJSONRequestBody) productCreateInputPort {
	return productCreateInputPort{
		params: p,
	}
}

func (p productCreateInputPort) Name() string {
	return p.params.Name
}

func (p productCreateInputPort) Price() int {
	return int(p.params.Price)
}
