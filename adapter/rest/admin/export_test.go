package admin

import "github.com/ryutah/virtual-ec/adapter/rest/admin/internal"

func NewServer(product *ProductEndpoint) *server {
	return &server{
		ProductEndpoint: product,
	}
}

type ProductSearchInputPort = productSearchInputPort

func NewProductSearchInputPort(p internal.ProductSearchParams) ProductSearchInputPort {
	return newProductSearchInputPort(p)
}

type ProductCreateInputPort = productCreateInputPort

func NewProductCreateInputPort(p internal.ProductCreateJSONRequestBody) ProductCreateInputPort {
	return newProductCreateInputPort(p)
}
