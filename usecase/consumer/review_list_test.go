package consumer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	. "github.com/ryutah/virtual-ec/usecase/consumer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReviewListOutputPort struct {
	mock.Mock
}

var _ ReviewListOutputPort = (*mockReviewListOutputPort)(nil)

func (m *mockReviewListOutputPort) onSuccess(reviewListSuccess interface{}) *mock.Call {
	return m.On("Success", reviewListSuccess)
}

func (m *mockReviewListOutputPort) onFailed(reviewListFailed interface{}) *mock.Call {
	return m.On("Failed", reviewListFailed)
}

func (m *mockReviewListOutputPort) Success(r ReviewListSuccess) {
	m.Called(r)
}

func (m *mockReviewListOutputPort) Failed(r ReviewListFailed) {
	m.Called(r)
}

func TestReviewList_List(t *testing.T) {
	type (
		mocks struct {
			repository_review_search_reviewSearchResult *repository.ReviewSearchResult
		}
		expected struct {
			args_repository_review_search_reviewQuery   repository.ReviewQuery
			args_reviewListOutputPort_reviewListSuccess ReviewListSuccess
		}
	)
	cases := []struct {
		name     string
		in       int
		mocks    mocks
		expected expected
	}{
		{
			name: "正常系",
			in:   2,
			mocks: mocks{
				repository_review_search_reviewSearchResult: &repository.ReviewSearchResult{
					Reviews: []*model.Review{
						model.ReCreateReview(1, 2, "user1", 5, "Good!"),
						model.ReCreateReview(10, 2, "user2", 3, "Bad!"),
					},
				},
			},
			expected: expected{
				args_repository_review_search_reviewQuery: repository.NewReviewQuery().WithProductID(2),
				args_reviewListOutputPort_reviewListSuccess: ReviewListSuccess{
					Reviews: []*ReviewListReviewDetail{
						{
							ID:       1,
							ReviewTo: 2,
							PostedBy: "user1",
							Rating:   5,
							Comment:  "Good!",
						},
						{
							ID:       10,
							ReviewTo: 2,
							PostedBy: "user2",
							Rating:   3,
							Comment:  "Bad!",
						},
					},
				}},
		},
		{
			name: "結果0件",
			in:   2,
			mocks: mocks{
				repository_review_search_reviewSearchResult: &repository.ReviewSearchResult{},
			},
			expected: expected{
				args_repository_review_search_reviewQuery: repository.NewReviewQuery().WithProductID(2),
				args_reviewListOutputPort_reviewListSuccess: ReviewListSuccess{
					Reviews: nil,
				}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			reviewRepo := new(mockReviewRepository)
			reviewRepo.
				onSearch(ctx, c.expected.args_repository_review_search_reviewQuery).
				Return(c.mocks.repository_review_search_reviewSearchResult, nil)
			output := new(mockReviewListOutputPort)
			output.onSuccess(c.expected.args_reviewListOutputPort_reviewListSuccess)

			review := NewReviewList(reviewRepo)
			success := review.List(ctx, c.in, output)

			reviewRepo.AssertExpectations(t)
			output.AssertExpectations(t)
			assert.True(t, success)
		})
	}
}

func TestReviewList_List_Failed(t *testing.T) {
	dummyError := errors.New("error")
	ctx := context.Background()

	reviewRepo := new(mockReviewRepository)
	reviewRepo.onSearch(ctx, mock.Anything).Return(nil, dummyError)
	output := new(mockReviewListOutputPort)
	output.onFailed(ReviewListFailed{
		Err: ReviewListErrorMessages.Failed(1),
	})

	review := NewReviewList(reviewRepo)
	success := review.List(ctx, 1, output)

	reviewRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}
