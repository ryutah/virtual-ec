package admin

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
)

type ProductEndpoint struct {
	usecase struct {
		searcher internal.ProductSearcher
		finder   internal.ProductFinder
	}
}

func NewProductEndpoint(seacher internal.ProductSearcher) *ProductEndpoint {
	return &ProductEndpoint{
		usecase: struct {
			searcher internal.ProductSearcher
			finder   internal.ProductFinder
		}{
			searcher: seacher,
		},
	}
}

// (GET /products)
func (p *ProductEndpoint) ProductSearch(w http.ResponseWriter, r *http.Request, params internal.ProductSearchParams) {
	out := new(productSearchOutputPort)
	_ = p.usecase.searcher.Search(r.Context(), newProductSearchInputPort(params), out)
	w.WriteHeader(out.status())
	render.JSON(w, r, out.payload())
}

// (GET /products/{product_id})
func (p *ProductEndpoint) ProductGet(w http.ResponseWriter, r *http.Request, productId int64) {
	panic("not implemented") // TODO: Implement
}

// (POST /products)
func (p *ProductEndpoint) ProductCreate(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}
