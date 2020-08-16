package admin

import "github.com/ryutah/virtual-ec/usecase/admin"

type productSearchInputPort struct {
	params ProductSearchParams
}

var _ admin.ProductSearchInputPort = productSearchInputPort{}

func newProductSearchInputPort(p ProductSearchParams) productSearchInputPort {
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
