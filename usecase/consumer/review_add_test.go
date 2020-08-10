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

func TestReviewAdd_Add(t *testing.T) {
	type (
		in struct {
			productID int
			req       ReviewAddRequest
		}
		mocks struct {
			repository_review_nextID model.ReviewID
		}
		expected struct {
			args_repository_review_nextID_productID model.ProductID
			args_repository_review_store_review     model.Review
			reviewAddResponse                       *ReviewAddResponse
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
				productID: 4,
				req: ReviewAddRequest{
					ReviewTo: 2,
					PostedBy: "user1",
					Rating:   3,
					Comment:  "Good!",
				},
			},
			mocks: mocks{
				repository_review_nextID: 1,
			},
			expected: expected{
				args_repository_review_nextID_productID: 4,
				args_repository_review_store_review: *model.NewReview(
					1, 2, "user1", 3, "Good!",
				),
				reviewAddResponse: &ReviewAddResponse{
					ID:       1,
					ReviewTo: 2,
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

			reviewRepo := new(mockReviewRepository)
			reviewRepo.On("NextID", ctx, c.expected.args_repository_review_nextID_productID).Return(c.mocks.repository_review_nextID, nil)
			reviewRepo.On("Store", ctx, c.expected.args_repository_review_store_review).Return(nil)

			review := NewReviewAdd(reviewRepo)
			got, err := review.Add(ctx, c.in.productID, c.in.req)

			reviewRepo.AssertExpectations(t)
			assert.Equal(t, c.expected.reviewAddResponse, got)
			assert.Nil(t, err)
		})
	}
}

func TestReviewAdd_Add_Failed_NextID(t *testing.T) {
	dummyError := errors.New("error")

	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything, mock.Anything).Return(model.ReviewID(0), dummyError)

	review := NewReviewAdd(reviewRepo)
	got, err := review.Add(context.Background(), 1, ReviewAddRequest{})

	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}

func TestReviewAdd_Add_Failed_Store(t *testing.T) {
	dummyError := errors.New("error")

	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything, mock.Anything).Return(model.ReviewID(1), nil)
	reviewRepo.On("Store", mock.Anything, mock.Anything).Return(dummyError)

	review := NewReviewAdd(reviewRepo)
	got, err := review.Add(context.Background(), 1, ReviewAddRequest{})

	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}
