package consumer

import (
	"context"
	"fmt"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/ryutah/virtual-ec/lib/xlog"
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
		Err string
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

func (r *ReviewList) List(ctx context.Context, productID int, out ReviewListOutputPort) (success bool) {
	result, err := r.repo.review.Search(
		ctx, repository.NewReviewQuery().WithProductID(model.ProductID(productID)),
	)
	if err != nil {
		xlog.Errorf(ctx, "failed to get review list: %+v", err)
		out.Failed(ReviewListFailed{
			Err: reviewListErrorMessages.failed(model.ProductID(productID)),
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
	out.Success(ReviewListSuccess{
		Reviews: responseReviews,
	})
	return true
}
