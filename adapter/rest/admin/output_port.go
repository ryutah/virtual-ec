package admin

import (
	"net/http"

	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
	"github.com/ryutah/virtual-ec/usecase/admin"
)

type outputPortBase struct {
	_status  int
	_payload interface{}
}

func (p outputPortBase) status() int {
	return p._status
}

func (p outputPortBase) payload() interface{} {
	return p._payload
}

type productSearchOutputPort struct {
	outputPortBase
}

var _ admin.ProductSearchOutputPort = (*productSearchOutputPort)(nil)

func (p *productSearchOutputPort) Success(success admin.ProductSearchSuccess) {
	products := make([]internal.Product, len(success.Products))
	for i, p := range success.Products {
		products[i] = internal.Product{
			Id:    int64(p.ID),
			Name:  p.Name,
			Price: int64(p.Price),
		}
	}
	p._payload = internal.ProductSearchSuccess{
		Products: products,
	}
	p._status = http.StatusOK
}

func (p *productSearchOutputPort) Failed(failed admin.ProductSearchFailed) {
	p._payload = internal.ServerError{
		Message: failed.Err,
	}
	p._status = http.StatusInternalServerError
}

type productFindOutputPort struct {
	outputPortBase
}

var _ admin.ProductFindOutputPort = (*productFindOutputPort)(nil)

func (p *productFindOutputPort) Success(success admin.ProductFindSuccess) {
	p._payload = internal.ProductGetSuccess{
		Id:    int64(success.ID),
		Name:  success.Name,
		Price: int64(success.Price),
	}
	p._status = http.StatusOK
}

func (p *productFindOutputPort) NotFound(failed admin.ProductFindFailed) {
	p._payload = internal.NotFound{
		Message: failed.Err,
	}
	p._status = http.StatusNotFound
}

func (p *productFindOutputPort) Failed(failed admin.ProductFindFailed) {
	p._payload = internal.NotFound{
		Message: failed.Err,
	}
	p._status = http.StatusInternalServerError
}

type productCreateOutputPort struct {
	outputPortBase
}

var _ admin.ProductCreateOutputPort = (*productCreateOutputPort)(nil)

func (p *productCreateOutputPort) Success(success admin.ProductCreateSuccess) {
	p._payload = internal.ProductCreateSuccess{
		Id:    int64(success.ID),
		Name:  success.Name,
		Price: int64(success.Price),
	}
	p._status = http.StatusCreated
}

func (p *productCreateOutputPort) Failed(_ admin.ProductCreateFailed) {
	panic("not implemented") // TODO: Implement
}
