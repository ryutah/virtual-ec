package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

var reviewErrMessages = struct {
	nextID func(error) string
	store  func(model.Review, error) string
	search func(repository.ReviewQuery, error) string
}{
	nextID: func(err error) string {
		return fmt.Sprintf("failed to allocates id: %v", err)
	},
	store: func(r model.Review, err error) string {
		return fmt.Sprintf("failed to store review(%v): %v", r, err)
	},
	search: func(q repository.ReviewQuery, err error) string {
		return fmt.Sprintf("failed to search(%v): %v", q, err)
	},
}

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
		return 0, errors.New(reviewErrMessages.nextID(err))
	}
	return model.ReviewID(keys[0].ID), nil
}

func (r *Review) Store(ctx context.Context, review model.Review) error {
	entity := reviewEntity{
		PostedBy: review.PostedBy(),
		Rating:   review.Rating(),
		Comment:  review.Comment(),
	}
	if _, err := r.client.Put(ctx, reviewKey(review.ReviewTo(), review.ID()), &entity); err != nil {
		return errors.New(reviewErrMessages.store(review, err))
	}
	return nil
}

func (r *Review) Search(ctx context.Context, query repository.ReviewQuery) (*repository.ReviewSearchResult, error) {
	q := datastore.NewQuery(kinds.review)
	if productID, ok := query.ProductID(); ok {
		q = q.Ancestor(productKey(productID))
	}

	entities := make([]*reviewEntity, 0)
	keys, err := r.client.GetAll(ctx, q, &entities)
	if err != nil {
		return nil, errors.New(reviewErrMessages.search(query, err))
	}

	var reviews []*model.Review
	for i, key := range keys {
		reviews = append(reviews, model.ReCreateReview(
			model.ReviewID(key.ID), model.ProductID(key.Parent.ID),
			entities[i].PostedBy, entities[i].Rating, entities[i].Comment,
		))
	}

	return &repository.ReviewSearchResult{
		Reviews: reviews,
	}, nil
}

func reviewKey(productID model.ProductID, reviewID model.ReviewID) *datastore.Key {
	return datastore.IDKey(kinds.review, int64(reviewID), productKey(productID))
}
