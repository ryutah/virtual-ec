package admin_test

import (
	"context"
	"errors"
	"testing"

	perrors "github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/usecase/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductFindOutputPort struct {
	mock.Mock
}

var _ ProductFindOutputPort = (*mockProductFindOutputPort)(nil)

func (m *mockProductFindOutputPort) onSuccess(productFindSuccess interface{}) *mock.Call {
	return m.On("Success", productFindSuccess)
}

func (m *mockProductFindOutputPort) onNotFound(productFindFailed interface{}) *mock.Call {
	return m.On("NotFound", productFindFailed)
}

func (m *mockProductFindOutputPort) onFailed(productFindFailed interface{}) *mock.Call {
	return m.On("Failed", productFindFailed)
}

func (m *mockProductFindOutputPort) Success(p ProductFindSuccess) {
	m.Called(p)
}

func (m *mockProductFindOutputPort) NotFound(p ProductFindFailed) {
	m.Called(p)
}

func (m *mockProductFindOutputPort) Failed(p ProductFindFailed) {
	m.Called(p)
}

func TestProductFind_Find(t *testing.T) {
	type (
		mocks struct {
			repository_product_get_product *model.Product
		}
		expected struct {
			args_repository_product_get_id                   model.ProductID
			args_product_find_output_port_productFindSuccess ProductFindSuccess
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
			in:   1,
			mocks: mocks{
				repository_product_get_product: model.NewProduct(
					1, "Product1", 100,
				),
			},
			expected: expected{
				args_repository_product_get_id: 1,
				args_product_find_output_port_productFindSuccess: ProductFindSuccess{
					ID:    1,
					Name:  "Product1",
					Price: 100,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			productRepo := new(mockProductRepository)
			productRepo.
				onGet(ctx, c.expected.args_repository_product_get_id).
				Return(c.mocks.repository_product_get_product, nil)
			output := new(mockProductFindOutputPort)
			output.onSuccess(c.expected.args_product_find_output_port_productFindSuccess)

			productFind := NewProductFind(output, productRepo)
			success := productFind.Find(ctx, c.in)

			productRepo.AssertExpectations(t)
			output.AssertExpectations(t)
			assert.True(t, success)
		})
	}
}

func TestProductFind_Find_Get_Filed(t *testing.T) {
	dummyError := errors.New("error")
	ctx := context.Background()

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)
	output := new(mockProductFindOutputPort)
	output.onFailed(ProductFindFailed{
		Err: errors.New(ProductFindFailedErrorMessages.Failed(1)),
	})

	productFind := NewProductFind(output, productRepo)
	success := productFind.Find(ctx, 1)

	productRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}

func TestProductFind_Find_Get_NotFound(t *testing.T) {
	dummyError := perrors.Wrap(domain.ErrNoSuchEntity, "dummy error")
	ctx := context.Background()

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)
	output := new(mockProductFindOutputPort)
	output.onNotFound(ProductFindFailed{
		Err: errors.New(ProductFindFailedErrorMessages.NotFound(1)),
	})

	productFind := NewProductFind(output, productRepo)
	success := productFind.Find(ctx, 1)

	productRepo.AssertExpectations(t)
	output.AssertExpectations(t)
	assert.False(t, success)
}
