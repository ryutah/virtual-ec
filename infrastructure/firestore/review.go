package firestore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type reviewEntity struct {
	PostedBy string
	Rating   int
	Comment  string
}

type Review struct {
	client Client
}

var _ repository.Review = (*Review)(nil)

func NewReview(c Client) *Review {
	return &Review{
		client: c,
	}
}

func (r *Review) NextID(ctx context.Context, productID model.ProductID) (model.ReviewID, error) {
	keys, err := r.client.AllocateIDs(ctx, []*datastore.Key{
		datastore.IncompleteKey(kinds.review, productKey(productID)),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return model.ReviewID(keys[0].ID), nil
}

func (r *Review) Store(ctx context.Context, review model.Review) error {
	entity := reviewEntity{
		PostedBy: review.PostedBy(),
		Rating:   review.Rating(),
		Comment:  review.Comment(),
	}
	_, err := r.client.Put(ctx, reviewKey(review.ReviewTo(), review.ID()), &entity)
	return errors.WithStack(err)
}

func (r *Review) Search(ctx context.Context, query repository.ReviewQuery) (*repository.ReviewSearchResult, error) {
	panic("not implement yet")
}

func reviewKey(productID model.ProductID, reviewID model.ReviewID) *datastore.Key {
	return datastore.IDKey(kinds.review, int64(reviewID), productKey(productID))
}