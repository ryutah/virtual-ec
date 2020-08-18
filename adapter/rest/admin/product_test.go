package admin_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func strPtr(s string) *string {
	return &s
}

func TestProduct_Search(t *testing.T) {
	type (
		mocks struct {
			call_productSearchOutputPort_productSearchSuccess admin.ProductSearchSuccess
		}
		expected struct {
			productSearch_search_input ProductSearchInputPort
			statusCode                 int
			response                   internal.ProductSearchSuccess
		}
	)
	cases := []struct {
		name     string
		in       url.Values
		mocks    mocks
		expected expected
	}{
		{
			name: "複数件取得",
			in: url.Values{
				"name": []string{"product"},
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
				response: internal.ProductSearchSuccess{
					Products: []internal.Product{
						{Id: 1, Name: "product1", Price: 1000},
						{Id: 2, Name: "product2", Price: 2000},
						{Id: 3, Name: "product3", Price: 3000},
					},
				},
			},
		},
		{
			name: "0件",
			in: url.Values{
				"name": []string{"product"},
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
				response: internal.ProductSearchSuccess{
					Products: []internal.Product{},
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

			handler := NewHandler(NewServer(NewProductEndpoint(searcher, new(mockProductFinder))))
			s := httptest.NewServer(handler)
			defer s.Close()

			req, _ := http.NewRequestWithContext(
				ctx, http.MethodGet, fmt.Sprintf("%s/api/products?%s", s.URL, c.in.Encode()), nil,
			)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer resp.Body.Close()

			var payload internal.ProductSearchSuccess
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

	handler := NewHandler(NewServer(NewProductEndpoint(searcher, new(mockProductFinder))))
	s := httptest.NewServer(handler)
	defer s.Close()

	req, _ := http.NewRequestWithContext(
		ctx, http.MethodGet, fmt.Sprintf("%s/api/products", s.URL), nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload internal.ServerError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	searcher.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, internal.ServerError{Message: dummeyErr.Error()}, payload)
}

func TestProduct_Get(t *testing.T) {
	type (
		mocks struct {
			call_productFindOutpuPort_productGetSucecss admin.ProductFindSuccess
		}
		expected struct {
			productFind_get_productID int
			statusCode                int
			response                  internal.ProductGetSuccess
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
				call_productFindOutpuPort_productGetSucecss: admin.ProductFindSuccess{
					ID:    1,
					Name:  "product1",
					Price: 1000,
				},
			},
			expected: expected{
				productFind_get_productID: 1,
				statusCode:                http.StatusOK,
				response: internal.ProductGetSuccess{
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

			handler := NewHandler(NewServer(NewProductEndpoint(new(mockProductSearcher), finder)))
			s := httptest.NewServer(handler)
			defer s.Close()

			req, _ := http.NewRequestWithContext(
				ctx, http.MethodGet, fmt.Sprintf("%s/api/products/%d", s.URL, c.in), nil,
			)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				assert.Fail(t, err.Error())
			}
			defer resp.Body.Close()

			var payload internal.ProductGetSuccess
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

	handler := NewHandler(NewServer(NewProductEndpoint(new(mockProductSearcher), finder)))
	s := httptest.NewServer(handler)
	defer s.Close()

	req, _ := http.NewRequestWithContext(
		ctx, http.MethodGet, fmt.Sprintf("%s/api/products/%d", s.URL, 1), nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload internal.NotFound
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	finder.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, internal.NotFound{Message: "not found"}, payload)
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

	handler := NewHandler(NewServer(NewProductEndpoint(new(mockProductSearcher), finder)))
	s := httptest.NewServer(handler)
	defer s.Close()

	req, _ := http.NewRequestWithContext(
		ctx, http.MethodGet, fmt.Sprintf("%s/api/products/%d", s.URL, 1), nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	defer resp.Body.Close()

	var payload internal.ServerError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		assert.Fail(t, err.Error())
	}

	finder.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, internal.ServerError{Message: "server error"}, payload)
}
