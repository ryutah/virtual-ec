package admin_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/usecase/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestProductCreator_Append(t *testing.T) {
	type (
		mocks struct {
			repository_product_nextID model.ProductID
		}
		expected struct {
			args_repository_product_store_product model.Product
			productAddResponse                    *ProductCreateResponse
			error                                 error
		}
	)
	cases := []struct {
		name     string
		in       ProductCreateRequest
		mocks    mocks
		expected expected
	}{
		{
			name: "指定した商品登録が正常に終了すること",
			in: ProductCreateRequest{
				Name:  "product1",
				Price: 1000,
			},
			mocks: mocks{
				repository_product_nextID: 1,
			},
			expected: expected{
				args_repository_product_store_product: *model.NewProduct(1, "product1", 1000),
				productAddResponse: &ProductCreateResponse{
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
			ctx := context.Background()

			productRepo := new(mockProductRepository)
			productRepo.onNextID(ctx).Return(c.mocks.repository_product_nextID, nil)
			productRepo.onStore(ctx, c.expected.args_repository_product_store_product).Return(nil)

			creator := NewProductCreate(productRepo)

			resp, err := creator.Create(ctx, c.in)
			assert.Equal(t, c.expected.productAddResponse, resp)
			assert.Equal(t, c.expected.error, err)
			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductCreator_Append_NextID_Failed(t *testing.T) {
	dummyError := errors.New("error")
	productRepo := new(mockProductRepository)
	productRepo.onNextID(mock.Anything).Return(model.ProductID(0), dummyError)

	creator := NewProductCreate(productRepo)

	resp, err := creator.Create(context.Background(), ProductCreateRequest{
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
	productRepo.onNextID(mock.Anything).Return(model.ProductID(1), nil)
	productRepo.onStore(mock.Anything, mock.Anything).Return(dummyError)

	creator := NewProductCreate(productRepo)

	resp, err := creator.Create(context.Background(), ProductCreateRequest{
		Name:  "test",
		Price: 100,
	})
	assert.Nil(t, resp)
	assert.Equal(t, dummyError, err)
	productRepo.AssertExpectations(t)
}
