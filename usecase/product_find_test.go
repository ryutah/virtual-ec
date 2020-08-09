package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductFind_Find(t *testing.T) {
	type (
		mocks struct {
			repository_product_get_product *model.Product
		}
		expected struct {
			args_repository_product_get_id model.ProductID
			response                       *ProductFindResponse
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
				response: &ProductFindResponse{
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

			productFind := NewProductFind(productRepo)
			got, err := productFind.Find(ctx, c.in)

			productRepo.AssertExpectations(t)
			assert.Equal(t, c.expected.response, got)
			assert.Nil(t, err)
		})
	}
}

func TestProductFind_Find_Get_Filed(t *testing.T) {
	dummyError := errors.New("error")
	ctx := context.Background()

	productRepo := new(mockProductRepository)
	productRepo.onGet(mock.Anything, mock.Anything).Return(nil, dummyError)

	productFind := NewProductFind(productRepo)
	got, err := productFind.Find(ctx, 1)

	productRepo.AssertExpectations(t)
	assert.Nil(t, got)
	assert.Equal(t, dummyError, err)
}
