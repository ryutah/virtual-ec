package repository

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
)

type ReviewQuery struct {
	productID *model.ProductID
}

func NewReviewQuery() ReviewQuery {
	return ReviewQuery{}
}

func (r ReviewQuery) WithProductID(id model.ProductID) ReviewQuery {
	r.productID = &id
	return r
}

func (r ReviewQuery) ProductID() (model.ProductID, bool) {
	if r.productID == nil {
		return 0, false
	}
	return *r.productID, true
}

type ReviewSearchResult struct {
	Reviews []*model.Review
}

type Review interface {
	NextID(context.Context, model.ProductID) (model.ReviewID, error)
	Search(context.Context, ReviewQuery) (*ReviewSearchResult, error)
	Store(context.Context, model.Review) error
}
