package usecase_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	. "github.com/ryutah/virtual-ec/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) NextID(ctx context.Context) (model.ProductID, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.ProductID), args.Error(1)
}

func (m *mockProductRepository) Store(ctx context.Context, p model.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

var _ repository.Product = (*mockProductRepository)(nil)

func TestProductCreator_Append(t *testing.T) {
	type (
		mocks struct {
			repository_product_nextID model.ProductID
		}
		expected struct {
			args_repository_product_store_product model.Product
			productAddResponse                    *ProductAddResponse
			error                                 error
		}
	)
	cases := []struct {
		name     string
		in       ProductAddRequest
		mocks    mocks
		expected expected
	}{
		{
			name: "指定した商品登録が正常に終了すること",
			in: ProductAddRequest{
				Name:  "product1",
				Price: 1000,
			},
			mocks: mocks{
				repository_product_nextID: 1,
			},
			expected: expected{
				args_repository_product_store_product: *model.NewProduct(1, "product1", 1000),
				productAddResponse: &ProductAddResponse{
					ID:    1,
					Name:  "product1",
					Price: 1000,
				},
				error: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			productRepo := new(mockProductRepository)
			productRepo.On("NextID", mock.Anything).Return(c.mocks.repository_product_nextID, nil)
			productRepo.On("Store", mock.Anything, c.expected.args_repository_product_store_product).Return(nil)

			creator := NewProductCreator(productRepo)

			resp, err := creator.Append(context.Background(), c.in)
			assert.Equal(t, c.expected.productAddResponse, resp)
			assert.Equal(t, c.expected.error, err)
			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductCreator_Append_NextID_Failed(t *testing.T) {
	dummyError := errors.New("error")
	productRepo := new(mockProductRepository)
	productRepo.On("NextID", mock.Anything).Return(model.ProductID(0), dummyError)

	creator := NewProductCreator(productRepo)

	resp, err := creator.Append(context.Background(), ProductAddRequest{
		Name:  "test",
		Price: 100,
	})
	assert.Nil(t, resp)
	assert.Equal(t, dummyError, err)
	productRepo.AssertExpectations(t)
}

func TestProductCreator_Append_Store_Failed(t *testing.T) {
	dummyError := errors.New("error")
	productRepo := new(mockProductRepository)
	productRepo.On("NextID", mock.Anything).Return(model.ProductID(1), nil)
	productRepo.On("Store", mock.Anything, mock.Anything).Return(dummyError)

	creator := NewProductCreator(productRepo)

	resp, err := creator.Append(context.Background(), ProductAddRequest{
		Name:  "test",
		Price: 100,
	})
	assert.Nil(t, resp)
	assert.Equal(t, dummyError, err)
	productRepo.AssertExpectations(t)
}
