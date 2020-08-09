package usecase

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type ReviewAddRequest struct {
	ReviewTo model.ProductID
	PostedBy string
	Rating   int
	Comment  string
}

type ReviewAddResponse struct {
	ID       model.ReviewID
	ReviewTo model.ProductID
	PostedBy string
	Rating   int
	Comment  string
}

type ReviewAdd struct {
	repo struct {
		review repository.Review
	}
}

func NewReviewAdd(reviewRepo repository.Review) *ReviewAdd {
	return &ReviewAdd{
		repo: struct{ review repository.Review }{
			review: reviewRepo,
		},
	}
}

func (r *ReviewAdd) Add(ctx context.Context, productID int, req ReviewAddRequest) (*ReviewAddResponse, error) {
	id, err := r.repo.review.NextID(ctx, model.ProductID(productID))
	if err != nil {
		return nil, err
	}
	review := model.NewReview(
		id, req.ReviewTo, req.PostedBy, req.Rating, req.Comment,
	)
	if err := r.repo.review.Store(ctx, *review); err != nil {
		return nil, err
	}

	return &ReviewAddResponse{
		ID:       review.ID(),
		ReviewTo: review.ReviewTo(),
		PostedBy: review.PostedBy(),
		Rating:   review.Rating(),
		Comment:  review.Comment(),
	}, nil
}
