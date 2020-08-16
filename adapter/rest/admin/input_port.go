package admin

import (
	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
	"github.com/ryutah/virtual-ec/usecase/admin"
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
