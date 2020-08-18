package admin_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	. "github.com/ryutah/virtual-ec/adapter/rest/admin"
	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
	"github.com/ryutah/virtual-ec/usecase/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductSearcher struct {
	internal.ProductSearcher
	mock.Mock
}

func (m *mockProductSearcher) onSearch(ctx, adminProductSearchInputPort, adminProductSearchOutputPort interface{}) *mock.Call {
	return m.On("Search", ctx, adminProductSearchInputPort, adminProductSearchOutputPort)
}

func (m *mockProductSearcher) Search(ctx context.Context, in admin.ProductSearchInputPort, out admin.ProductSearchOutputPort) bool {
	args := m.Called(ctx, in, out)
	return args.Bool(0)
}

type mockProductFinder struct {
	internal.ProductFinder
	mock.Mock
}

func (m *mockProductFinder) onFind(ctx, productID, productFindOutputPort interface{}) *mock.Call {
	return m.On("Find", ctx, productID, productFindOutputPort)
}

func (m *mockProductFinder) Find(ctx context.Context, productID int, out admin.ProductFindOutputPort) bool {
	args := m.Called(ctx, productID, out)
	return args.Bool(0)
}

func TestProduct_Search(t *testing.T) {
	type (
		mocks struct {
			call_productSearchOutputPort_productSearchSuccess admin.ProductSearchSuccess
		}
		expected struct {
			productSearch_search_input ProductSearchInputPort
			statusCode                 int
			response                   ProductSearchSuccess
		}
	)
	cases := []struct {
		name     string
		in       *ProductSearchParams
		mocks    mocks
		expected expected
	}{
		{
			name: "複数件取得",
			in: &ProductSearchParams{
				Name: strPtr("product"),
			},
			mocks: mocks{
				call_productSearchOutputPort_productSearchSuccess: admin.ProductSearchSuccess{
					Products: []*admin.ProductSearchProductDetail{
						{ID: 1, Name: "product1", Price: 1000},
						{ID: 2, Name: "product2", Price: 2000},
						{ID: 3, Name: "product3", Price: 3000},
					},
				},
			},
			expected: expected{
				productSearch_search_input: NewProductSearchInputPort(internal.ProductSearchParams{
					Name: strPtr("product"),
				}),
				statusCode: http.StatusOK,
				response: ProductSearchSuccess{
					Products: []Product{
						{Id: 1, Name: "product1", Price: 1000},
						{Id: 2, Name: "product2", Price: 2000},
						{Id: 3, Name: "product3", Price: 3000},
					},
				},
			},
		},
		{
			name: "検索条件未指定",
			in:   &ProductSearchParams{},
			mocks: mocks{
				call_productSearchOutputPort_productSearchSuccess: admin.ProductSearchSuccess{
					Products: []*admin.ProductSearchProductDetail{
						{ID: 1, Name: "product1", Price: 1000},
						{ID: 2, Name: "product2", Price: 2000},
						{ID: 3, Name: "product3", Price: 3000},
					},
				},
			},
			expected: expected{
				productSearch_search_input: NewProductSearchInputPort(internal.ProductSearchParams{}),
				statusCode:                 http.StatusOK,
				response: ProductSearchSuccess{
					Products: []Product{
						{Id: 1, Name: "product1", Price: 1000},
						{Id: 2, Name: "product2", Price: 2000},
						{Id: 3, Name: "product3", Price: 3000},
					},
				},
			},
		},
		{
			name: "0件",
			in: &ProductSearchParams{
				Name: strPtr("product"),
			},
			mocks: mocks{
				call_productSearchOutputPort_productSearchSuccess: admin.ProductSearchSuccess{
					Products: nil,
				},
			},
			expected: expected{
				productSearch_search_input: NewProductSearchInputPort(internal.ProductSearchParams{
					Name: strPtr("product"),
				}),
				statusCode: http.StatusOK,
				response: ProductSearchSuccess{
					Products: []Product{},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			// define usecase mock and set call behavior
			searcher := new(mockProductSearcher)
			searcher.
				onSearch(mock.Anything, c.expected.productSearch_search_input, mock.Anything).
				Return(true).
				Run(func(args mock.Arguments) {
					out := args.Get(2).(admin.ProductSearchOutputPort)
					out.Success(c.mocks.call_productSearchOutputPort_productSearchSuccess)
				})

			client, finish := startTestServerAndNewClient(withProductSearcher(searcher))
			defer finish()

			resp, err := client.ProductSearch(ctx, c.in)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer resp.Body.Close()

			var payload ProductSearchSuccess
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				assert.Fail(t, err.Error())
			}

			searcher.AssertExpectations(t)
			assert.Equal(t, c.expected.statusCode, resp.StatusCode)
			assert.Equal(t, c.expected.response, payload)
		})
	}
}

func TestProduct_Search_Failed(t *testing.T) {
	dummeyErr := errors.New("error")
	ctx := context.Background()

	// define usecase mock and set call behavior
	searcher := new(mockProductSearcher)
	searcher.
		onSearch(mock.Anything, mock.Anything, mock.Anything).
		Return(false).
		Run(func(args mock.Arguments) {
			out := args.Get(2).(admin.ProductSearchOutputPort)
			out.Failed(admin.ProductSearchFailed{
				Err: dummeyErr.Error(),
			})
		})

	client, finish := startTestServerAndNewClient(withProductSearcher(searcher))
	defer finish()

	resp, err := client.ProductSearch(ctx, new(ProductSearchParams))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload ServerError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	searcher.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, ServerError{Message: dummeyErr.Error()}, payload)
}

func TestProduct_Get(t *testing.T) {
	type (
		mocks struct {
			call_productFindOutpuPort_productGetSucecss admin.ProductFindSuccess
		}
		expected struct {
			productFind_get_productID int
			statusCode                int
			response                  ProductGetSuccess
		}
	)
	cases := []struct {
		name     string
		in       int64
		mocks    mocks
		expected expected
	}{
		{
			name: "正常系",
			in:   1,
			mocks: mocks{
				call_productFindOutpuPort_productGetSucecss: admin.ProductFindSuccess{
					ID:    1,
					Name:  "product1",
					Price: 1000,
				},
			},
			expected: expected{
				productFind_get_productID: 1,
				statusCode:                http.StatusOK,
				response: ProductGetSuccess{
					Id:    1,
					Name:  "product1",
					Price: 1000,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			// define usecase mock and set call behavior
			finder := new(mockProductFinder)
			finder.
				onFind(mock.Anything, c.expected.productFind_get_productID, mock.Anything).
				Return(true).
				Run(func(args mock.Arguments) {
					out := args.Get(2).(admin.ProductFindOutputPort)
					out.Success(c.mocks.call_productFindOutpuPort_productGetSucecss)
				})

			client, finish := startTestServerAndNewClient(withProductFinder(finder))
			defer finish()

			resp, err := client.ProductGet(ctx, c.in)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer resp.Body.Close()

			var payload ProductGetSuccess
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				assert.Fail(t, err.Error())
			}

			finder.AssertExpectations(t)
			assert.Equal(t, c.expected.statusCode, resp.StatusCode)
			assert.Equal(t, c.expected.response, payload)
		})
	}
}

func TestProduct_Get_NotFound(t *testing.T) {
	ctx := context.Background()

	// define usecase mock and set call behavior
	finder := new(mockProductFinder)
	finder.
		onFind(mock.Anything, mock.Anything, mock.Anything).
		Return(true).
		Run(func(args mock.Arguments) {
			out := args.Get(2).(admin.ProductFindOutputPort)
			out.NotFound(admin.ProductFindFailed{
				Err: "not found",
			})
		})

	client, finish := startTestServerAndNewClient(withProductFinder(finder))
	defer finish()

	resp, err := client.ProductGet(ctx, 1)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload NotFound
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	finder.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, NotFound{Message: "not found"}, payload)
}

func TestProduct_Get_Failed(t *testing.T) {
	ctx := context.Background()

	// define usecase mock and set call behavior
	finder := new(mockProductFinder)
	finder.
		onFind(mock.Anything, mock.Anything, mock.Anything).
		Return(true).
		Run(func(args mock.Arguments) {
			out := args.Get(2).(admin.ProductFindOutputPort)
			out.Failed(admin.ProductFindFailed{
				Err: "server error",
			})
		})

	client, finish := startTestServerAndNewClient(withProductFinder(finder))
	defer finish()

	resp, err := client.ProductGet(ctx, 1)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload ServerError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	finder.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, ServerError{Message: "server error"}, payload)
}
