package firestore_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/ryutah/virtual-ec/domain/model"
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
	err := review.Store(ctx, *model.NewReview(1, 2, "user1", 3, "comments!"))

	client.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestReview_Store_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	ctx := context.Background()

	client := new(mockClient)
	client.onPut(mock.Anything, mock.Anything, mock.Anything).Return(nil, dummyErr)

	review := NewReview(client)
	err := review.Store(ctx, *model.NewReview(1, 2, "user1", 3, "comments!"))

	client.AssertExpectations(t)
	assert.EqualError(t, err, dummyErr.Error())
}
