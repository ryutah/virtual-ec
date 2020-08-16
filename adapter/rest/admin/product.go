package admin

import (
	"net/http"

	"github.com/go-chi/render"
)

type ProductEndpoint struct {
	usecase struct {
		searcher ProductSearcher
		finder   ProductFinder
	}
}

func NewProductEndpoint(seacher ProductSearcher) *ProductEndpoint {
	return &ProductEndpoint{
		usecase: struct {
			searcher ProductSearcher
			finder   ProductFinder
		}{
			searcher: seacher,
		},
	}
}

// (GET /products)
func (p *ProductEndpoint) ProductSearch(w http.ResponseWriter, r *http.Request, params ProductSearchParams) {
	out := new(productSearchOutputPort)
	_ = p.usecase.searcher.Search(r.Context(), newProductSearchInputPort(params), out)
	w.WriteHeader(out.status())
	render.JSON(w, r, out.payload())
}

// (GET /products/{product_id})
func (p *ProductEndpoint) ProductGet(w http.ResponseWriter, r *http.Request, productId int64) {
	panic("not implemented") // TODO: Implement
}
