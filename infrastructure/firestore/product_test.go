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

type mockClient struct {
	Client
	mock.Mock
}

var _ Client = (*mockClient)(nil)

func (m *mockClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	args := m.Called(ctx, keys)
	retKeys, ok := args.Get(0).([]*datastore.Key)
	if !ok {
		retKeys = nil
	}
	return retKeys, args.Error(1)
}

func (m *mockClient) Put(ctx context.Context, key *datastore.Key, v interface{}) (*datastore.Key, error) {
	args := m.Called(ctx, key, v)
	retKey, ok := args.Get(0).(*datastore.Key)
	if !ok {
		retKey = nil
	}
	return retKey, args.Error(1)
}

func TestProduct_NextID(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient)
	client.
		On("AllocateIDs", ctx, []*datastore.Key{
			datastore.IncompleteKey(Kinds.Product, nil),
		}).
		Return([]*datastore.Key{
			datastore.IDKey(Kinds.Product, 1, nil),
		}, nil)

	product := NewProduct(client)
	got, err := product.NextID(ctx)

	client.AssertExpectations(t)
	assert.Equal(t, model.ProductID(1), got)
	assert.Nil(t, err)
}

func TestProduct_NextID_Failed(t *testing.T) {
	dummyErr := errors.New("error")

	client := new(mockClient)
	client.On("AllocateIDs", mock.Anything, mock.Anything).Return(nil, dummyErr)

	product := NewProduct(client)
	got, err := product.NextID(context.Background())

	client.AssertExpectations(t)
	assert.Zero(t, got)
	assert.EqualError(t, err, dummyErr.Error())
}

func TestProduct_Store(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient)
	client.
		On("Put", ctx, datastore.IDKey(Kinds.Product, 1, nil), &ProductEntity{
			Name:  "product",
			Price: 100,
		}).
		Return(datastore.IDKey(Kinds.Product, 1, nil), nil)

	product := NewProduct(client)
	err := product.Store(ctx, *model.NewProduct(1, "product", 100))

	client.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestProduct_Store_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	client := new(mockClient)
	client.On("Put", mock.Anything, mock.Anything, mock.Anything).Return(nil, dummyErr)

	product := NewProduct(client)
	err := product.Store(context.Background(), *model.NewProduct(1, "product", 100))

	client.AssertExpectations(t)
	assert.EqualError(t, err, dummyErr.Error())
}
