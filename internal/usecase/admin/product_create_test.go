package admin_test

import (
	"testing"

	"errors"

	"github.com/ryutah/virtual-ec/internal/domain/model"
	. "github.com/ryutah/virtual-ec/internal/usecase/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type mockProductCreateOutputPort struct {
	mock.Mock
}

var _ ProductCreateOutputPort = (*mockProductCreateOutputPort)(nil)

func (m *mockProductCreateOutputPort) onSuccess(productCreateSuccess interface{}) *mock.Call {
	return m.On("Success", productCreateSuccess)
}

func (m *mockProductCreateOutputPort) onFailed(productCreateFailed interface{}) *mock.Call {
	return m.On("Failed", productCreateFailed)
}

func (m *mockProductCreateOutputPort) Success(p ProductCreateSuccess) {
	m.Called(p)
}

func (m *mockProductCreateOutputPort) Failed(p ProductCreateFailed) {
	m.Called(p)
}

type productCreateImputPort struct {
	name  string
	price int
}

var _ ProductCreateInputPort = (*productCreateImputPort)(nil)

func (p productCreateImputPort) Name() string {
	return p.name
}

func (p productCreateImputPort) Price() int {
	return p.price
}

func TestProductCreator_Append(t *testing.T) {
	type (
		mocks struct {
			repository_product_nextID model.ProductID
		}
		expected struct {
			args_repository_product_store_product             model.Product
			args_productCreateOutputPort_productCreateSuccess ProductCreateSuccess
			productAddResponse                                *ProductCreateSuccess
			error                                             error
		}
	)
	cases := []struct {
		name     string
		in       ProductCreateInputPort
		mocks    mocks
		expected expected
	}{
		{
			name: "指定した商品登録が正常に終了すること",
			in: productCreateImputPort{
				name:  "product1",
				price: 1000,
			},
			mocks: mocks{
				repository_product_nextID: 1,
			},
			expected: expected{
				args_repository_product_store_product: *model.ReCreateProduct(1, "product1", 1000),
				args_productCreateOutputPort_productCreateSuccess: ProductCreateSuccess{
					ID:    1,
					Name:  "product1",
					Price: 1000,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			productRepo := new(mockProductRepository)
			productRepo.onNextID(ctx).Return(c.mocks.repository_product_nextID, nil)
			productRepo.onStore(ctx, c.expected.args_repository_product_store_product).Return(nil)
			output := new(mockProductCreateOutputPort)
			output.onSuccess(c.expected.args_productCreateOutputPort_productCreateSuccess)

			creator := NewProductCreate(productRepo)
			success := creator.Create(ctx, c.in, output)

			productRepo.AssertExpectations(t)
			output.AssertExpectations(t)
			assert.True(t, success)
		})
	}
}

func TestProductCreator_Append_NextID_Failed(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onNextID(mock.Anything).Return(model.ProductID(0), dummyError)
	output := new(mockProductCreateOutputPort)
	output.onFailed(ProductCreateFailed{
		Err: ProductCreateErrroMessages.Failed(),
	})

	creator := NewProductCreate(productRepo)
	success := creator.Create(context.Background(), productCreateImputPort{
		name:  "test",
		price: 100,
	}, output)

	productRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}

func TestProductCreator_Append_Store_Failed(t *testing.T) {
	dummyError := errors.New("error")

	productRepo := new(mockProductRepository)
	productRepo.onNextID(mock.Anything).Return(model.ProductID(100), nil)
	productRepo.onStore(mock.Anything, mock.Anything).Return(dummyError)
	output := new(mockProductCreateOutputPort)
	output.onFailed(ProductCreateFailed{
		Err: ProductCreateErrroMessages.Failed(),
	})

	creator := NewProductCreate(productRepo)
	success := creator.Create(context.Background(), productCreateImputPort{
		name:  "test",
		price: 100,
	}, output)

	productRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}
