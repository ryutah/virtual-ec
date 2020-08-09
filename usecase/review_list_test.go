package usecase_test

import (
	"context"
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	. "github.com/ryutah/virtual-ec/usecase"
	"github.com/stretchr/testify/assert"
)

func TestReviewList_List(t *testing.T) {
	type (
		mocks struct {
			repository_review_search_reviewSearchResult *repository.ReviewSearchResult
		}
		expected struct {
			args_repository_review_search_reviewQuery repository.ReviewQuery
			reviewListResponse                        *ReviewListResponse
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
						model.NewReview(1, 2, "user1", 5, "Good!"),
						model.NewReview(10, 2, "user2", 3, "Bad!"),
					},
				},
			},
			expected: expected{
				args_repository_review_search_reviewQuery: repository.NewReviewQuery().WithProductID(2),
				reviewListResponse: &ReviewListResponse{
					Reviews: []*ReviewListResponseReview{
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
				reviewListResponse: &ReviewListResponse{
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

			review := NewReviewList(reviewRepo)
			got, err := review.List(ctx, c.in)

			assert.Equal(t, c.expected.reviewListResponse, got)
			assert.Nil(t, err)
		})
	}
}
