//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductFind -as internal.ProductFinder -o internal/product_finder_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductSearch -as internal.ProductSearcher -o internal/product_seacher_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductCreate -as internal.ProductCreator -o internal/product_creator_gen.go
//go:generate oapi-codegen -generate chi-server,types -package internal -o internal/openapi_gen.go ../../../documents/admin/openapi.yaml

package admin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
)

func NewHandler(s internal.ServerInterface) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.StripSlashes)

	mux.Route("/api", func(r chi.Router) {
		internal.HandlerFromMux(s, r)
	})
	return mux
}

type Server struct {
	*ProductEndpoint
}

var _ internal.ServerInterface = (*Server)(nil)

func NewServer(product *ProductEndpoint) *Server {
	return &Server{
		ProductEndpoint: product,
	}
}
