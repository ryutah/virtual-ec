package firestore_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	perrors "github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/internal/domain"
	"github.com/ryutah/virtual-ec/internal/domain/model"
	"github.com/ryutah/virtual-ec/internal/domain/repository"
	. "github.com/ryutah/virtual-ec/internal/infrastructure/firestore"
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
	err := product.Store(ctx, *model.ReCreateProduct(1, "product", 100))

	client.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestProduct_Store_Failed(t *testing.T) {
	dummyErr := errors.New("error")
	client := new(mockClient)
	client.onPut(mock.Anything, mock.Anything, mock.Anything).Return(nil, dummyErr)

	modelProduct := model.ReCreateProduct(1, "product", 100)
	product := NewProduct(client)
	err := product.Store(context.Background(), *modelProduct)

	client.AssertExpectations(t)
	assert.EqualError(t, err, ProductErrMessages.Store(*modelProduct, dummyErr))
}

func TestProduct_Get(t *testing.T) {
	ctx := context.Background()

	client := new(mockClient).withStore(productStore{
		{
			k: ProductKey(1),
			v: &ProductEntity{
				Name:  "Product1",
				Price: 100,
			},
		},
	})
	client.onGet(ctx, ProductKey(1), new(ProductEntity)).Return(nil)

	product := NewProduct(client)
	got, err := product.Get(ctx, 1)

	client.AssertExpectations(t)
	assert.Equal(t, model.ReCreateProduct(1, "Product1", 100), got)
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

func TestProduct_Search(t *testing.T) {
	type expected struct {
		args_client_run_datastoreQuery *datastore.Query
		productSearchResult            *repository.ProductSearchResult
	}
	cases := []struct {
		name     string
		store    productStore
		in       repository.ProductQuery
		expected expected
	}{
		{
			name: "複数件取得",
			store: productStore{
				{
					k: ProductKey(1),
					v: &ProductEntity{Name: "product1", Price: 1},
				},
				{
					k: ProductKey(2),
					v: &ProductEntity{Name: "product2", Price: 2},
				},
				{
					k: ProductKey(3),
					v: &ProductEntity{Name: "product3", Price: 3},
				},
			},
			in: repository.NewProductQuery().WithName("product"),
			expected: expected{
				args_client_run_datastoreQuery: datastore.NewQuery(Kinds.Product).Filter("Name=", "product"),
				productSearchResult: &repository.ProductSearchResult{
					Products: []*model.Product{
						model.ReCreateProduct(1, "product1", 1),
						model.ReCreateProduct(2, "product2", 2),
						model.ReCreateProduct(3, "product3", 3),
					},
				},
			},
		},
		{
			name:  "結果0件",
			store: productStore{},
			in:    repository.NewProductQuery().WithName("product"),
			expected: expected{
				args_client_run_datastoreQuery: datastore.NewQuery(Kinds.Product).Filter("Name=", "product"),
				productSearchResult: &repository.ProductSearchResult{
					Products: nil,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			it := new(mockIterator).withStore(c.store)
			it.onNext(new(ProductEntity)).Times(len(c.store) + 1)
			client := new(mockClient)
			client.onRun(ctx, c.expected.args_client_run_datastoreQuery).Return(it)

			product := NewProduct(client)
			got, err := product.Search(ctx, c.in)

			it.AssertExpectations(t)
			client.AssertExpectations(t)
			assert.Equal(t, c.expected.productSearchResult, got)
			assert.Nil(t, err)
		})
	}
}

func TestProduct_Search_Failed(t *testing.T) {
	ctx := context.Background()

	it := new(mockIterator).withStore(productStore{})
	it.onNext(new(ProductEntity)).Return(nil, errors.New("error")).Times(1)
	client := new(mockClient)
	client.onRun(ctx, mock.Anything).Return(it)

	product := NewProduct(client)
	got, err := product.Search(ctx, repository.NewProductQuery())

	it.AssertExpectations(t)
	client.AssertExpectations(t)
	assert.Nil(t, got)
	assert.EqualError(t, err, ProductErrMessages.Search(errors.New("error")))
}
