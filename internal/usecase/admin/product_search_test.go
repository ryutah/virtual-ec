package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ryutah/virtual-ec/internal/domain/model"
	"github.com/ryutah/virtual-ec/internal/domain/repository"
	. "github.com/ryutah/virtual-ec/internal/usecase/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductSearchOutputPort struct {
	mock.Mock
}

var _ ProductSearchOutputPort = (*mockProductSearchOutputPort)(nil)

func (m *mockProductSearchOutputPort) onSuccess(productSearchSuccess interface{}) *mock.Call {
	return m.On("Success", productSearchSuccess)
}

func (m *mockProductSearchOutputPort) onFailed(productSearchFailed interface{}) *mock.Call {
	return m.On("Failed", productSearchFailed)
}

func (m *mockProductSearchOutputPort) Success(p ProductSearchSuccess) {
	m.Called(p)
}

func (m *mockProductSearchOutputPort) Failed(p ProductSearchFailed) {
	m.Called(p)
}

type productSearchInutPort struct {
	name string
}

var _ ProductSearchInputPort = (*productSearchInutPort)(nil)

func (p productSearchInutPort) Name() string {
	return p.name
}

func TestProductSearch_Search(t *testing.T) {
	type (
		mocks struct {
			repository_product_search_produtSearchResult *repository.ProductSearchResult
		}
		expected struct {
			args_repository_product_search_productQuery      repository.ProductQuery
			args_productSearchOutputPort_produtSearchSuccess ProductSearchSuccess
		}
	)
	cases := []struct {
		name     string
		in       ProductSearchInputPort
		mocks    mocks
		expected expected
	}{
		{
			name: "複数件取得",
			in: productSearchInutPort{
				name: "product",
			},
			mocks: mocks{
				repository_product_search_produtSearchResult: &repository.ProductSearchResult{
					Products: []*model.Product{
						model.ReCreateProduct(1, "product1", 100),
						model.ReCreateProduct(2, "product2", 200),
						model.ReCreateProduct(3, "product3", 300),
					},
				},
			},
			expected: expected{
				args_repository_product_search_productQuery: repository.NewProductQuery().WithName("product"),
				args_productSearchOutputPort_produtSearchSuccess: ProductSearchSuccess{
					Products: []*ProductSearchProductDetail{
						{ID: 1, Name: "product1", Price: 100},
						{ID: 2, Name: "product2", Price: 200},
						{ID: 3, Name: "product3", Price: 300},
					},
				},
			},
		},
		{
			name: "検索条件未指定",
			in:   productSearchInutPort{},
			mocks: mocks{
				repository_product_search_produtSearchResult: &repository.ProductSearchResult{
					Products: []*model.Product{
						model.ReCreateProduct(1, "product1", 100),
						model.ReCreateProduct(2, "product2", 200),
						model.ReCreateProduct(3, "product3", 300),
					},
				},
			},
			expected: expected{
				args_repository_product_search_productQuery: repository.NewProductQuery(),
				args_productSearchOutputPort_produtSearchSuccess: ProductSearchSuccess{
					Products: []*ProductSearchProductDetail{
						{ID: 1, Name: "product1", Price: 100},
						{ID: 2, Name: "product2", Price: 200},
						{ID: 3, Name: "product3", Price: 300},
					},
				},
			},
		},
		{
			name: "結果0件",
			in:   productSearchInutPort{},
			mocks: mocks{
				repository_product_search_produtSearchResult: &repository.ProductSearchResult{
					Products: nil,
				},
			},
			expected: expected{
				args_repository_product_search_productQuery: repository.NewProductQuery(),
				args_productSearchOutputPort_produtSearchSuccess: ProductSearchSuccess{
					Products: nil,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			productRepo := new(mockProductRepository)
			productRepo.
				onSearch(ctx, c.expected.args_repository_product_search_productQuery).
				Return(c.mocks.repository_product_search_produtSearchResult, nil)
			output := new(mockProductSearchOutputPort)
			output.onSuccess(c.expected.args_productSearchOutputPort_produtSearchSuccess)

			product := NewProductSearch(productRepo)
			success := product.Search(ctx, c.in, output)

			productRepo.AssertExpectations(t)
			output.AssertExpectations(t)
			assert.True(t, success)
		})
	}
}

func TestProductSearch_Search_Failed(t *testing.T) {
	ctx := context.Background()

	productRepo := new(mockProductRepository)
	productRepo.onSearch(ctx, mock.Anything).Return(mock.Anything, errors.New("error"))
	output := new(mockProductSearchOutputPort)
	output.onFailed(ProductSearchFailed{
		Err: ProductSearchErrorMessages.Failed(),
	})

	product := NewProductSearch(productRepo)
	success := product.Search(ctx, productSearchInutPort{}, output)

	productRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}
