//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductFind -as admin.ProductFinder -o product_finder_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductSearch -as admin.ProductSearcher -o product_seacher_gen.go
//go:generate interfacer -for github.com/ryutah/virtual-ec/usecase/admin.ProductCreate -as admin.ProductCreator -o product_creator_gen.go

package admin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewHandler(s ServerInterface) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.StripSlashes)

	mux.Route("/api", func(r chi.Router) {
		HandlerFromMux(s, r)
	})
	return mux
}

type Server struct {
	*ProductEndpoint
}

var _ ServerInterface = (*Server)(nil)

func NewServer(product *ProductEndpoint) *Server {
	return &Server{
		ProductEndpoint: product,
	}
}
