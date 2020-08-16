package admin

import "github.com/ryutah/virtual-ec/adapter/rest/admin/internal"

type ProductSearchInputPort = productSearchInputPort

func NewProductSearchInputPort(p internal.ProductSearchParams) ProductSearchInputPort {
	return newProductSearchInputPort(p)
}
