package admin

type ProductSearchInputPort = productSearchInputPort

func NewProductSearchInputPort(p ProductSearchParams) ProductSearchInputPort {
	return newProductSearchInputPort(p)
}
