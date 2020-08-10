package consumer

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type ReviewPostRequest struct {
	PostedBy string
	Rating   int
	Comment  string
}

type ReviewPostResponse struct {
	ID       model.ReviewID
	ReviewTo model.ProductID
	PostedBy string
	Rating   int
	Comment  string
}

type ReviewPost struct {
	repo struct {
		review  repository.Review
		product repository.Product
	}
}

func NewReviewPost(reviewRepo repository.Review, productRepo repository.Product) *ReviewPost {
	return &ReviewPost{
		repo: struct {
			review  repository.Review
			product repository.Product
		}{
			review:  reviewRepo,
			product: productRepo,
		},
	}
}

func (r *ReviewPost) Post(ctx context.Context, productID int, req ReviewPostRequest) (*ReviewPostResponse, error) {
	product, err := r.repo.product.Get(ctx, model.ProductID(productID))
	if err != nil {
		return nil, err
	}
	id, err := r.repo.review.NextID(ctx, model.ProductID(productID))
	if err != nil {
		return nil, err
	}

	review := product.NewReview(id)
	review.Write(req.PostedBy, req.Rating, req.Comment)

	if err := r.repo.review.Store(ctx, *review); err != nil {
		return nil, err
	}

	return &ReviewPostResponse{
		ID:       review.ID(),
		ReviewTo: review.ReviewTo(),
		PostedBy: review.PostedBy(),
		Rating:   review.Rating(),
		Comment:  review.Comment(),
	}, nil
}
