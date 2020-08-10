package firestore_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	. "github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReview_NextID(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient)
	client.
		onAllocateIDs(ctx, []*datastore.Key{
			datastore.IncompleteKey(
				Kinds.Review, ProductKey(2),
			),
		}).
		Return([]*datastore.Key{ReviewKey(2, 1)}, nil)

	review := NewReview(client)
	got, err := review.NextID(ctx, 2)

	client.AssertExpectations(t)
	assert.Equal(t, model.ReviewID(1), got)
	assert.Nil(t, err)
}

func TestReview_NextID_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	ctx := context.Background()

	client := new(mockClient)
	client.onAllocateIDs(mock.Anything, mock.Anything).Return(nil, dummyErr)

	review := NewReview(client)
	got, err := review.NextID(ctx, 1)

	client.AssertExpectations(t)
	assert.Zero(t, got)
	assert.EqualError(t, err, dummyErr.Error())
}

func TestReview_Store(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient)
	client.onPut(ctx, ReviewKey(2, 1), &ReviewEntity{
		PostedBy: "user1",
		Rating:   3,
		Comment:  "comments!",
	}).Return(ReviewKey(2, 1), nil)

	review := NewReview(client)
	err := review.Store(ctx, *model.ReCreateReview(1, 2, "user1", 3, "comments!"))

	client.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestReview_Store_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	ctx := context.Background()

	client := new(mockClient)
	client.onPut(mock.Anything, mock.Anything, mock.Anything).Return(nil, dummyErr)

	review := NewReview(client)
	err := review.Store(ctx, *model.ReCreateReview(1, 2, "user1", 3, "comments!"))

	client.AssertExpectations(t)
	assert.EqualError(t, err, dummyErr.Error())
}

func TestReview_Search(t *testing.T) {
	type (
		in    struct{}
		mocks struct {
			repository_review_getAll_datastoreKeys []*datastore.Key
		}
		expected struct {
			args_repository_review_getAll_query *datastore.Query
			args_repository_review_getAll_value *[]*ReviewEntity
			reviewSearchResult                  *repository.ReviewSearchResult
		}
	)
	ptrOf := func(r []*ReviewEntity) *[]*ReviewEntity {
		return &r
	}

	cases := []struct {
		name     string
		in       repository.ReviewQuery
		store    reviewStore
		mocks    mocks
		expected expected
	}{
		{
			name: "正常系",
			in:   repository.NewReviewQuery().WithProductID(1),
			store: reviewStore{
				{
					key: ReviewKey(1, 1),
					val: &ReviewEntity{
						PostedBy: "user1",
						Rating:   1,
						Comment:  "comment1",
					},
				},
				{
					key: ReviewKey(1, 2),
					val: &ReviewEntity{
						PostedBy: "user2",
						Rating:   2,
						Comment:  "comment2",
					},
				},
				{
					key: ReviewKey(1, 3),
					val: &ReviewEntity{
						PostedBy: "user3",
						Rating:   3,
						Comment:  "comment3",
					},
				},
			},
			mocks: mocks{
				repository_review_getAll_datastoreKeys: []*datastore.Key{
					ReviewKey(1, 1),
					ReviewKey(1, 2),
					ReviewKey(1, 3),
				},
			},
			expected: expected{
				args_repository_review_getAll_query: datastore.NewQuery(Kinds.Review).Ancestor(ProductKey(1)),
				args_repository_review_getAll_value: ptrOf(make([]*ReviewEntity, 0)),
				reviewSearchResult: &repository.ReviewSearchResult{
					Reviews: []*model.Review{
						model.ReCreateReview(1, 1, "user1", 1, "comment1"),
						model.ReCreateReview(2, 1, "user2", 2, "comment2"),
						model.ReCreateReview(3, 1, "user3", 3, "comment3"),
					},
				},
			},
		},
		{
			name:  "結果0件",
			in:    repository.NewReviewQuery().WithProductID(1),
			store: nil,
			mocks: mocks{
				repository_review_getAll_datastoreKeys: nil,
			},
			expected: expected{
				args_repository_review_getAll_query: datastore.NewQuery(Kinds.Review).Ancestor(ProductKey(1)),
				args_repository_review_getAll_value: ptrOf(make([]*ReviewEntity, 0)),
				reviewSearchResult: &repository.ReviewSearchResult{
					Reviews: nil,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			client := new(mockClient).withStore(c.store)
			client.
				onGetAll(ctx, c.expected.args_repository_review_getAll_query, c.expected.args_repository_review_getAll_value).
				Return(c.mocks.repository_review_getAll_datastoreKeys, nil)

			review := NewReview(client)
			got, err := review.Search(ctx, c.in)

			client.AssertExpectations(t)
			assert.Equal(t, c.expected.reviewSearchResult, got)
			assert.Nil(t, err)
		})
	}
}

func TestReview_Search_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	ctx := context.Background()

	client := new(mockClient)
	client.onGetAll(ctx, mock.Anything, mock.Anything).Return(nil, dummyErr)

	review := NewReview(client)
	got, err := review.Search(ctx, repository.NewReviewQuery())

	client.AssertExpectations(t)
	assert.Nil(t, got)
	assert.EqualError(t, err, dummyErr.Error())
}
