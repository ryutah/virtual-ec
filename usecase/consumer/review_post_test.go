package consumer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/usecase/consumer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReviewPost_Post(t *testing.T) {
	type (
		in struct {
			productID int
			req       ReviewPostRequest
		}
		mocks struct {
			repository_product_get_product *model.Product
			repository_review_nextID       model.ReviewID
		}
		expected struct {
			args_repository_product_get_productID   model.ProductID
			args_repository_review_nextID_productID model.ProductID
			args_repository_review_store_review     model.Review
			reviewAddResponse                       *ReviewPostResponse
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
				req: ReviewPostRequest{
					PostedBy: "user1",
					Rating:   3,
					Comment:  "Good!",
				},
			},
			mocks: mocks{
				repository_product_get_product: model.NewProduct(1, "product", 100),
				repository_review_nextID:       2,
			},
			expected: expected{
				args_repository_product_get_productID:   1,
				args_repository_review_nextID_productID: 1,
				args_repository_review_store_review: *model.ReCreateReview(
					2, 1, "user1", 3, "Good!",
				),
				reviewAddResponse: &ReviewPostResponse{
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
			reviewRepo.On("NextID", ctx, c.expected.args_repository_review_nextID_productID).Return(c.mocks.repository_review_nextID, nil)
			reviewRepo.On("Store", ctx, c.expected.args_repository_review_store_review).Return(nil)

			review := NewReviewPost(reviewRepo, productRepo)
			got, err := review.Post(ctx, c.in.productID, c.in.req)

			productRepo.AssertExpectations(t)
			reviewRepo.AssertExpectations(t)
			assert.Equal(t, c.expected.reviewAddResponse, got)
			assert.Nil(t, err)
		})
	}
}

func TestReviewPost_Post_Failed_ProductGet(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)
	reviewRepo := new(mockReviewRepository)

	review := NewReviewPost(reviewRepo, productRepo)
	got, err := review.Post(context.Background(), 1, ReviewPostRequest{})

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}

func TestReviewPost_Post_Failed_NextID(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(model.NewProduct(1, "product", 100), nil)
	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything, mock.Anything).Return(model.ReviewID(0), dummyError)

	review := NewReviewPost(reviewRepo, productRepo)
	got, err := review.Post(context.Background(), 1, ReviewPostRequest{})

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}

func TestReviewPost_Post_Failed_Store(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(model.NewProduct(1, "product", 100), nil)
	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything, mock.Anything).Return(model.ReviewID(1), nil)
	reviewRepo.On("Store", mock.Anything, mock.Anything).Return(dummyError)

	review := NewReviewPost(reviewRepo, productRepo)
	got, err := review.Post(context.Background(), 1, ReviewPostRequest{})

	productRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}
