package admin_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/internal/adapter/rest/admin"
	"github.com/ryutah/virtual-ec/internal/adapter/rest/admin/internal"
	"github.com/stretchr/testify/assert"
)

func TestProductSearchInputPort(t *testing.T) {
	type expected struct {
		name string
	}
	cases := []struct {
		name     string
		in       internal.ProductSearchParams
		expected expected
	}{
		{
			name: "正常文字列",
			in: internal.ProductSearchParams{
				Name: strPtr("product"),
			},
			expected: expected{
				name: "product",
			},
		},
		{
			name: "文字列null",
			in: internal.ProductSearchParams{
				Name: nil,
			},
			expected: expected{
				name: "",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			port := NewProductSearchInputPort(c.in)
			assert.Equal(t, c.expected.name, port.Name())
		})
	}
}

func TestProductCreateInputPort(t *testing.T) {
	type expected struct {
		name  string
		price int
	}
	cases := []struct {
		name     string
		in       internal.ProductCreateJSONRequestBody
		expected expected
	}{
		{
			name: "正常系",
			in: internal.ProductCreateJSONRequestBody{
				Name:  "product",
				Price: 100,
			},
			expected: expected{
				name:  "product",
				price: 100,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			port := NewProductCreateInputPort(c.in)
			assert.Equal(t, c.expected.name, port.Name())
			assert.Equal(t, c.expected.price, port.Price())
		})
	}
}
