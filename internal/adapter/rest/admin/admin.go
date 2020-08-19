//go:generate interfacer -for github.com/ryutah/virtual-ec/internal/usecase/admin.ProductFind -as internal.ProductFinder -o internal/product_finder_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/internal/usecase/admin.ProductSearch -as internal.ProductSearcher -o internal/product_seacher_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/internal/usecase/admin.ProductCreate -as internal.ProductCreator -o internal/product_creator_gen.go
//go:generate oapi-codegen -generate chi-server,types -package internal -o internal/openapi_gen.go ../../../../api/admin/openapi.yaml

package admin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ryutah/virtual-ec/internal/adapter/rest/admin/internal"
)

func NewHandler(product *ProductEndpoint) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.StripSlashes)

	mux.Route("/api", func(r chi.Router) {
		internal.HandlerFromMux(newServer(product), r)
	})
	return mux
}

type server struct {
	*ProductEndpoint
}

var _ internal.ServerInterface = (*server)(nil)

func newServer(product *ProductEndpoint) *server {
	return &server{
		ProductEndpoint: product,
	}
}
