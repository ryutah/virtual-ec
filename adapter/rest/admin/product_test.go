package admin_test

import (
	"context"
	"encoding/json"
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

			handler := NewHandler(NewServer(NewProductEndpoint(searcher)))
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
			assert.Equal(t, c.expected.response, payload)
		})
	}
}
