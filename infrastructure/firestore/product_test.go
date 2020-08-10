package firestore_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	perrors "github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
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
	assert.EqualError(t, err, ProductErrMessages.NextID(dummyErr))
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

	modelProduct := model.NewProduct(1, "product", 100)
	product := NewProduct(client)
	err := product.Store(context.Background(), *modelProduct)

	client.AssertExpectations(t)
	assert.EqualError(t, err, ProductErrMessages.Store(*modelProduct, dummyErr))
}

func TestProduct_Get(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient).withStore(productStore{
		{
			key: ProductKey(1),
			val: &ProductEntity{
				Name:  "Product1",
				Price: 100,
			},
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
	type (
		mocks struct {
			client_get_error error
		}
	)
	cases := []struct {
		name     string
		in       model.ProductID
		mocks    mocks
		expected error
	}{
		{
			name: "指定したkeyに該当するエンティティが存在しない",
			in:   1,
			mocks: mocks{
				client_get_error: datastore.ErrNoSuchEntity,
			},
			expected: domain.ErrNoSuchEntity,
		},
		{
			name: "不明なエラーが発生",
			in:   1,
			mocks: mocks{
				client_get_error: errors.New("some error"),
			},
			expected: errors.New(ProductErrMessages.Get(1, errors.New("some error"))),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			client := new(mockClient)
			client.onGet(mock.Anything, mock.Anything, mock.Anything).Return(c.mocks.client_get_error)

			product := NewProduct(client)
			got, err := product.Get(ctx, c.in)

			client.AssertExpectations(t)
			assert.Nil(t, got)
			assert.EqualError(t, perrors.Cause(err), c.expected.Error())
		})
	}
}
