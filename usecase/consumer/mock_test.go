package consumer_test

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
	"github.com/stretchr/testify/mock"
)

type mockProductRepository struct {
	mock.Mock
}

var _ repository.Product = (*mockProductRepository)(nil)

func (m *mockProductRepository) onNextID(ctx interface{}) *mock.Call {
	return m.On("NextID", ctx)
}

func (m *mockProductRepository) onStore(ctx, modelProduct interface{}) *mock.Call {
	return m.On("Store", ctx, modelProduct)
}

func (m *mockProductRepository) onGet(ctx, id interface{}) *mock.Call {
	return m.On("Get", ctx, id)
}

func (m *mockProductRepository) NextID(ctx context.Context) (model.ProductID, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.ProductID), args.Error(1)
}

func (m *mockProductRepository) Store(ctx context.Context, p model.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockProductRepository) Get(ctx context.Context, id model.ProductID) (*model.Product, error) {
	args := m.Called(ctx, id)
	product, _ := args.Get(0).(*model.Product)
	return product, args.Error(1)
}

type mockReviewRepository struct {
	mock.Mock
}

var _ repository.Review = (*mockReviewRepository)(nil)

func (m *mockReviewRepository) onNextID(ctx, modelProductID interface{}) *mock.Call {
	return m.On("NextID", ctx, modelProductID)
}

func (m *mockReviewRepository) onStore(ctx, modelReview interface{}) *mock.Call {
	return m.On("Store", ctx, modelReview)
}

func (m *mockReviewRepository) onSearch(ctx, repositoryReviewQuery interface{}) *mock.Call {
	return m.On("Search", ctx, repositoryReviewQuery)
}

func (m *mockReviewRepository) NextID(ctx context.Context, pid model.ProductID) (model.ReviewID, error) {
	args := m.Called(ctx, pid)
	return args.Get(0).(model.ReviewID), args.Error(1)
}

func (m *mockReviewRepository) Store(ctx context.Context, review model.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *mockReviewRepository) Search(ctx context.Context, query repository.ReviewQuery) (*repository.ReviewSearchResult, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).(*repository.ReviewSearchResult)
	return result, args.Error(1)
}
