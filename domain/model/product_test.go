package model_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	type (
		in struct {
			id    ProductID
			name  string
			price int
		}
		expected struct {
			id    ProductID
			name  string
			price int
		}
	)
	cases := []struct {
		name     string
		in       in
		expected expected
	}{
		{
			name: "test1",
			in: in{
				id:    1,
				name:  "product1",
				price: 1000,
			},
			expected: expected{
				id:    1,
				name:  "product1",
				price: 1000,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NewProduct(c.in.id, c.in.name, c.in.price)
			assert.Equal(t, c.expected.id, got.ID())
			assert.Equal(t, c.expected.name, got.Name())
			assert.Equal(t, c.expected.price, got.Price())
		})
	}
}
