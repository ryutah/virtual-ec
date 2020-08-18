package admin

import "github.com/ryutah/virtual-ec/adapter/rest/admin/internal"

type ProductSearchInputPort = productSearchInputPort

func NewProductSearchInputPort(p internal.ProductSearchParams) ProductSearchInputPort {
	return newProductSearchInputPort(p)
}

func NewServer(product *ProductEndpoint) *server {
	return &server{
		ProductEndpoint: product,
	}
}
