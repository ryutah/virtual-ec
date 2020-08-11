package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

var reviewListErrorMessages = struct {
	failed func(model.ProductID) string
}{
	failed: func(id model.ProductID) string {
		return fmt.Sprintf("Product(%v)のレビューの取得に失敗しました", id)
	},
}

type ReviewListOutputPort interface {
	Success(ReviewListSuccess)
	Failed(ReviewListFailed)
}

type (
	ReviewListSuccess struct {
		Reviews []*ReviewListReviewDetail
	}

	ReviewListReviewDetail struct {
		ID       int
		ReviewTo int
		PostedBy string
		Rating   int
		Comment  string
	}

	ReviewListFailed struct {
		Err error
	}
)

type ReviewList struct {
	output ReviewListOutputPort
	repo   struct {
		review repository.Review
	}
}

func NewReviewList(output ReviewListOutputPort, reviewRepo repository.Review) *ReviewList {
	return &ReviewList{
		output: output,
		repo: struct{ review repository.Review }{
			review: reviewRepo,
		},
	}
}

func (r *ReviewList) List(ctx context.Context, productID int) (success bool) {
	result, err := r.repo.review.Search(
		ctx, repository.NewReviewQuery().WithProductID(model.ProductID(productID)),
	)
	if err != nil {
		r.output.Failed(ReviewListFailed{
			Err: errors.New(reviewListErrorMessages.failed(model.ProductID(productID))),
		})
		return false
	}

	var responseReviews []*ReviewListReviewDetail
	for _, r := range result.Reviews {
		responseReviews = append(responseReviews, &ReviewListReviewDetail{
			ID:       int(r.ID()),
			ReviewTo: int(r.ReviewTo()),
			PostedBy: r.PostedBy(),
			Rating:   r.Rating(),
			Comment:  r.Comment(),
		})
	}
	r.output.Success(ReviewListSuccess{
		Reviews: responseReviews,
	})
	return true
}
