package admin

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/ryutah/virtual-ec/internal/adapter/rest/admin/internal"
)

type ProductEndpoint struct {
	usecase struct {
		searcher internal.ProductSearcher
		finder   internal.ProductFinder
		creator  internal.ProductCreator
	}
}

func NewProductEndpoint(seacher internal.ProductSearcher, finder internal.ProductFinder, creator internal.ProductCreator) *ProductEndpoint {
	return &ProductEndpoint{
		usecase: struct {
			searcher internal.ProductSearcher
			finder   internal.ProductFinder
			creator  internal.ProductCreator
		}{
			searcher: seacher,
			finder:   finder,
			creator:  creator,
		},
	}
}

// (GET /products)
func (p *ProductEndpoint) ProductSearch(w http.ResponseWriter, r *http.Request, params internal.ProductSearchParams) {
	out := new(productSearchOutputPort)
	_ = p.usecase.searcher.Search(r.Context(), newProductSearchInputPort(params), out)
	render.Status(r, out.status())
	render.JSON(w, r, out.payload())
}

// (GET /products/{product_id})
func (p *ProductEndpoint) ProductGet(w http.ResponseWriter, r *http.Request, productId int64) {
	out := new(productFindOutputPort)
	_ = p.usecase.finder.Find(r.Context(), int(productId), out)
	render.Status(r, out.status())
	render.JSON(w, r, out.payload())
}

// (POST /products)
func (p *ProductEndpoint) ProductCreate(w http.ResponseWriter, r *http.Request) {
	var payload internal.ProductCreateJSONRequestBody
	_ = render.DecodeJSON(r.Body, &payload)

	out := new(productCreateOutputPort)
	_ = p.usecase.creator.Create(r.Context(), newProductCreateInputPort(payload), out)

	render.Status(r, out.status())
	render.JSON(w, r, out.payload())
}
