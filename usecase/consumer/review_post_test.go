package consumer_test

import (
	"context"
	"errors"
	"testing"

	perrors "github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/usecase/consumer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReviewPostOutputPort struct {
	mock.Mock
}

var _ ReviewPostOutputPort = (*mockReviewPostOutputPort)(nil)

func (m *mockReviewPostOutputPort) onSuccess(reviewPostSuccess interface{}) *mock.Call {
	return m.On("Success", reviewPostSuccess)
}

func (m *mockReviewPostOutputPort) onProductNotFound(reviewPostFailed interface{}) *mock.Call {
	return m.On("ProductNotFound", reviewPostFailed)
}

func (m *mockReviewPostOutputPort) onFailed(reviewPostFailed interface{}) *mock.Call {
	return m.On("Failed", reviewPostFailed)
}

func (m *mockReviewPostOutputPort) Success(r ReviewPostSuccess) {
	m.Called(r)
}

func (m *mockReviewPostOutputPort) ProductNotFound(r ReviewPostFailed) {
	m.Called(r)
}

func (m *mockReviewPostOutputPort) Failed(r ReviewPostFailed) {
	m.Called(r)
}

type reviewPostInputPort struct {
	postedBy string
	rating   int
	comment  string
}

func (r reviewPostInputPort) PostedBy() string {
	return r.postedBy
}

func (r reviewPostInputPort) Rating() int {
	return r.rating
}

func (r reviewPostInputPort) Comment() string {
	return r.comment
}

func TestReviewPost_Post(t *testing.T) {
	type (
		in struct {
			productID int
			input     ReviewPostInputPort
		}
		mocks struct {
			repository_product_get_product *model.Product
			repository_review_nextID       model.ReviewID
		}
		expected struct {
			args_repository_product_get_productID       model.ProductID
			args_repository_review_nextID_productID     model.ProductID
			args_repository_review_store_review         model.Review
			args_reviewPostOutputPort_reviewPostSuccess ReviewPostSuccess
		}
	)
	cases := []struct {
		name     string
		in       in
		mocks    mocks
		expected expected
	}{
		{
			name: "正常系",
			in: in{
				productID: 1,
				input: reviewPostInputPort{
					postedBy: "user1",
					rating:   3,
					comment:  "Good!",
				},
			},
			mocks: mocks{
				repository_product_get_product: model.ReCreateProduct(1, "product", 100),
				repository_review_nextID:       2,
			},
			expected: expected{
				args_repository_product_get_productID:   1,
				args_repository_review_nextID_productID: 1,
				args_repository_review_store_review: *model.ReCreateReview(
					2, 1, "user1", 3, "Good!",
				),
				args_reviewPostOutputPort_reviewPostSuccess: ReviewPostSuccess{
					ID:       2,
					ReviewTo: 1,
					PostedBy: "user1",
					Rating:   3,
					Comment:  "Good!",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			productRepo := new(mockProductRepository)
			productRepo.onGet(ctx, c.expected.args_repository_product_get_productID).Return(c.mocks.repository_product_get_product, nil)
			reviewRepo := new(mockReviewRepository)
			reviewRepo.onNextID(ctx, c.expected.args_repository_review_nextID_productID).Return(c.mocks.repository_review_nextID, nil)
			reviewRepo.onStore(ctx, c.expected.args_repository_review_store_review).Return(nil)
			output := new(mockReviewPostOutputPort)
			output.onSuccess(c.expected.args_reviewPostOutputPort_reviewPostSuccess)

			review := NewReviewPost(reviewRepo, productRepo)
			success := review.Post(ctx, c.in.productID, c.in.input, output)

			productRepo.AssertExpectations(t)
			reviewRepo.AssertExpectations(t)
			output.AssertExpectations(t)
			assert.True(t, success)
		})
	}
}

func TestReviewPost_Post_Failed_ProductGet(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)
	reviewRepo := new(mockReviewRepository)
	output := new(mockReviewPostOutputPort)
	output.onFailed(ReviewPostFailed{
		Err: ReviewPostErrorMessages.Failed(),
	})

	review := NewReviewPost(reviewRepo, productRepo)
	success := review.Post(context.Background(), 1, reviewPostInputPort{}, output)

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}

func TestReviewPost_Post_Failed_ProductGet_ProductNotFound(t *testing.T) {
	dummyError := perrors.Wrap(domain.ErrNoSuchEntity, "error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)
	reviewRepo := new(mockReviewRepository)
	output := new(mockReviewPostOutputPort)
	output.onProductNotFound(ReviewPostFailed{
		Err: ReviewPostErrorMessages.ProductNotFound(1),
	})

	review := NewReviewPost(reviewRepo, productRepo)
	success := review.Post(context.Background(), 1, reviewPostInputPort{}, output)

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}

func TestReviewPost_Post_Failed_NextID(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(model.ReCreateProduct(1, "product", 100), nil)
	reviewRepo := new(mockReviewRepository)
	reviewRepo.onNextID(mock.Anything, mock.Anything).Return(model.ReviewID(0), dummyError)
	output := new(mockReviewPostOutputPort)
	output.onFailed(ReviewPostFailed{
		Err: ReviewPostErrorMessages.Failed(),
	})

	review := NewReviewPost(reviewRepo, productRepo)
	success := review.Post(context.Background(), 1, reviewPostInputPort{}, output)

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}

func TestReviewPost_Post_Failed_Store(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(model.ReCreateProduct(1, "product", 100), nil)
	reviewRepo := new(mockReviewRepository)
	reviewRepo.onNextID(mock.Anything, mock.Anything).Return(model.ReviewID(1), nil)
	reviewRepo.onStore(mock.Anything, mock.Anything).Return(dummyError)
	output := new(mockReviewPostOutputPort)
	output.onFailed(ReviewPostFailed{
		Err: ReviewPostErrorMessages.Failed(),
	})

	review := NewReviewPost(reviewRepo, productRepo)
	success := review.Post(context.Background(), 1, reviewPostInputPort{}, output)

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}
