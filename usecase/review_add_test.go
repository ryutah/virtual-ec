package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	. "github.com/ryutah/virtual-ec/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReviewRepository struct {
	mock.Mock
}

var _ repository.Review = (*mockReviewRepository)(nil)

func (m *mockReviewRepository) NextID(ctx context.Context) (model.ReviewID, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.ReviewID), args.Error(1)
}

func (m *mockReviewRepository) Store(ctx context.Context, review model.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func TestReviewAdd_Add(t *testing.T) {
	type (
		mocks struct {
			repository_review_nextID model.ReviewID
		}
		expected struct {
			args_repository_review_store_review model.Review
			reviewAddResponse                   *ReviewAddResponse
		}
	)
	cases := []struct {
		name     string
		in       ReviewAddRequest
		mocks    mocks
		expected expected
	}{
		{
			name: "正常系",
			in: ReviewAddRequest{
				ReviewTo: 2,
				PostedBy: "user1",
				Rating:   3,
				Comment:  "Good!",
			},
			mocks: mocks{
				repository_review_nextID: 1,
			},
			expected: expected{
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
			reviewRepo.On("NextID", ctx).Return(c.mocks.repository_review_nextID, nil)
			reviewRepo.On("Store", ctx, c.expected.args_repository_review_store_review).Return(nil)

			review := NewReviewAdd(reviewRepo)
			got, err := review.Add(ctx, c.in)

			reviewRepo.AssertExpectations(t)
			assert.Equal(t, c.expected.reviewAddResponse, got)
			assert.Nil(t, err)
		})
	}
}

func TestReviewAdd_Add_Failed_NextID(t *testing.T) {
	dummyError := errors.New("error")

	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything).Return(model.ReviewID(0), dummyError)

	review := NewReviewAdd(reviewRepo)
	got, err := review.Add(context.Background(), ReviewAddRequest{})

	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}

func TestReviewAdd_Add_Failed_Store(t *testing.T) {
	dummyError := errors.New("error")

	reviewRepo := new(mockReviewRepository)
	reviewRepo.On("NextID", mock.Anything).Return(model.ReviewID(1), nil)
	reviewRepo.On("Store", mock.Anything, mock.Anything).Return(dummyError)

	review := NewReviewAdd(reviewRepo)
	got, err := review.Add(context.Background(), ReviewAddRequest{})

	reviewRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}
