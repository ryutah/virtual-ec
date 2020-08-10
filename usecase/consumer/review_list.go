package consumer

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type (
	ReviewListResponse struct {
		Reviews []*ReviewListResponseReview
	}

	ReviewListResponseReview struct {
		ID       int
		ReviewTo int
		PostedBy string
		Rating   int
		Comment  string
	}
)

type ReviewList struct {
	repo struct {
		review repository.Review
	}
}

func NewReviewList(reviewRepo repository.Review) *ReviewList {
	return &ReviewList{
		repo: struct{ review repository.Review }{
			review: reviewRepo,
		},
	}
}

func (r *ReviewList) List(ctx context.Context, productID int) (*ReviewListResponse, error) {
	result, _ := r.repo.review.Search(
		ctx, repository.NewReviewQuery().WithProductID(model.ProductID(productID)),
	)

	var responseReviews []*ReviewListResponseReview
	for _, r := range result.Reviews {
		responseReviews = append(responseReviews, &ReviewListResponseReview{
			ID:       int(r.ID()),
			ReviewTo: int(r.ReviewTo()),
			PostedBy: r.PostedBy(),
			Rating:   r.Rating(),
			Comment:  r.Comment(),
		})
	}

	return &ReviewListResponse{
		Reviews: responseReviews,
	}, nil
}
