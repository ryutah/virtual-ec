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

func TestProduct_NextID(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient)
	client.
		onAllocateIDs(ctx, []*datastore.Key{
			datastore.IncompleteKey(Kinds.Product, nil),
		}).
		Return([]*datastore.Key{
			ProductKey(1),
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
	client.onAllocateIDs(mock.Anything, mock.Anything).Return(nil, dummyErr)

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
		onPut(ctx, ProductKey(1), &ProductEntity{
			Name:  "product",
			Price: 100,
		}).
		Return(ProductKey(1), nil)

	product := NewProduct(client)
	err := product.Store(ctx, *model.NewProduct(1, "product", 100))

	client.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestProduct_Store_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	client := new(mockClient)
	client.onPut(mock.Anything, mock.Anything, mock.Anything).Return(nil, dummyErr)

	product := NewProduct(client)
	err := product.Store(context.Background(), *model.NewProduct(1, "product", 100))

	client.AssertExpectations(t)
	assert.EqualError(t, err, dummyErr.Error())
}

func TestProduct_Get(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient).withStore(mockClientStore{
		*ProductKey(1): &ProductEntity{
			Name:  "Product1",
			Price: 100,
		},
	})
	client.onGet(ctx, ProductKey(1), new(ProductEntity)).Return(nil)

	product := NewProduct(client)
	got, err := product.Get(ctx, 1)

	client.AssertExpectations(t)
	assert.Equal(t, model.NewProduct(1, "Product1", 100), got)
	assert.Nil(t, err)
}

func TestProduct_Get_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	ctx := context.Background()

	client := new(mockClient)
	client.onGet(mock.Anything, mock.Anything, mock.Anything).Return(dummyErr)

	product := NewProduct(client)
	got, err := product.Get(ctx, 1)

	client.AssertExpectations(t)
	assert.Nil(t, got)
	assert.EqualError(t, err, dummyErr.Error())
}
