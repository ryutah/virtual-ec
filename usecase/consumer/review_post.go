package consumer

import (
	"context"
	"errors"
	"fmt"

	perrors "github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

var reviewPostErrorMessages = struct {
	productNotFound func(model.ProductID) string
	failed          func() string
}{
	productNotFound: func(id model.ProductID) string {
		return fmt.Sprintf("Product(%v)は存在しません", id)
	},
	failed: func() string { return "Reviewの投稿に失敗しました" },
}

type (
	ReviewPostOutputPort interface {
		Success(ReviewPostSuccess)
		ProductNotFound(ReviewPostFailed)
		Failed(ReviewPostFailed)
	}

	ReviewPostInputPort interface {
		PostedBy() string
		Rating() int
		Comment() string
	}
)

type (
	ReviewPostSuccess struct {
		ID       model.ReviewID
		ReviewTo model.ProductID
		PostedBy string
		Rating   int
		Comment  string
	}

	ReviewPostFailed struct {
		Err error
	}
)

type ReviewPost struct {
	output ReviewPostOutputPort
	repo   struct {
		review  repository.Review
		product repository.Product
	}
}

func NewReviewPost(output ReviewPostOutputPort, reviewRepo repository.Review, productRepo repository.Product) *ReviewPost {
	return &ReviewPost{
		output: output,
		repo: struct {
			review  repository.Review
			product repository.Product
		}{
			review:  reviewRepo,
			product: productRepo,
		},
	}
}

func (r *ReviewPost) Post(ctx context.Context, productID int, input ReviewPostInputPort) (success bool) {
	product, err := r.repo.product.Get(ctx, model.ProductID(productID))
	if err != nil {
		return r.handleGetProductError(model.ProductID(productID), err)
	}
	id, err := r.repo.review.NextID(ctx, model.ProductID(productID))
	if err != nil {
		return r.handleError()
	}

	review := product.NewReview(id)
	review.Write(input.PostedBy(), input.Rating(), input.Comment())

	if err := r.repo.review.Store(ctx, *review); err != nil {
		return r.handleError()
	}

	r.output.Success(ReviewPostSuccess{
		ID:       review.ID(),
		ReviewTo: review.ReviewTo(),
		PostedBy: review.PostedBy(),
		Rating:   review.Rating(),
		Comment:  review.Comment(),
	})
	return true
}

func (r *ReviewPost) handleGetProductError(id model.ProductID, err error) bool {
	if perrors.Is(err, domain.ErrNoSuchEntity) {
		r.output.ProductNotFound(ReviewPostFailed{
			Err: errors.New(reviewPostErrorMessages.productNotFound(id)),
		})
	} else {
		r.output.Failed(ReviewPostFailed{
			Err: errors.New(reviewPostErrorMessages.failed()),
		})
	}
	return false
}

func (r *ReviewPost) handleError() bool {
	r.output.Failed(ReviewPostFailed{
		Err: errors.New(reviewPostErrorMessages.failed()),
	})
	return false
}
