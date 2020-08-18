package admin

import (
	"net/http"

	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
	"github.com/ryutah/virtual-ec/usecase/admin"
)

type resultGetter interface {
	status() int
	payload() interface{}
}

type productSearchOutputPort struct {
	statuses int
	payloads interface{}
}

var (
	_ admin.ProductSearchOutputPort = (*productSearchOutputPort)(nil)
	_ resultGetter                  = (*productSearchOutputPort)(nil)
)

func (p *productSearchOutputPort) Success(success admin.ProductSearchSuccess) {
	products := make([]internal.Product, len(success.Products))
	for i, p := range success.Products {
		products[i] = internal.Product{
			Id:    int64(p.ID),
			Name:  p.Name,
			Price: int64(p.Price),
		}
	}
	p.payloads = internal.ProductSearchSuccess{
		Products: products,
	}
	p.statuses = http.StatusOK
}

func (p *productSearchOutputPort) Failed(failed admin.ProductSearchFailed) {
	p.payloads = internal.ServerError{
		Message: failed.Err,
	}
	p.statuses = http.StatusInternalServerError
}

func (p *productSearchOutputPort) status() int {
	return p.statuses
}

func (p *productSearchOutputPort) payload() interface{} {
	return p.payloads
}

type productFindOutputPort struct {
	statuses int
	payloads interface{}
}

var (
	_ admin.ProductFindOutputPort = (*productFindOutputPort)(nil)
	_ resultGetter                = (*productFindOutputPort)(nil)
)

func (p *productFindOutputPort) Success(success admin.ProductFindSuccess) {
	p.payloads = internal.ProductGetSuccess{
		Id:    int64(success.ID),
		Name:  success.Name,
		Price: int64(success.Price),
	}
	p.statuses = http.StatusOK
}

func (p *productFindOutputPort) NotFound(failed admin.ProductFindFailed) {
	p.payloads = internal.NotFound{
		Message: failed.Err,
	}
	p.statuses = http.StatusNotFound
}

func (p *productFindOutputPort) Failed(failed admin.ProductFindFailed) {
	p.payloads = internal.NotFound{
		Message: failed.Err,
	}
	p.statuses = http.StatusInternalServerError
}

func (p *productFindOutputPort) status() int {
	return p.statuses
}

func (p *productFindOutputPort) payload() interface{} {
	return p.payloads
}
