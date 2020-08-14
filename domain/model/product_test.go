package model_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	. "github.com/ryutah/virtual-ec/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestProduct_NewProduct(t *testing.T) {
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
			name: "正常系",
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
			got, err := NewProduct(c.in.id, c.in.name, c.in.price)

			assert.Equal(t, c.expected.id, got.ID())
			assert.Equal(t, c.expected.name, got.Name())
			assert.Equal(t, c.expected.price, got.Price())
			assert.Nil(t, err)
		})
	}
}

func TestProduct_NewProduct_Failed(t *testing.T) {
	type in struct {
		id    ProductID
		name  string
		price int
	}
	cases := []struct {
		name string
		in   in
	}{
		{
			name: "name空欄",
			in: in{
				id:    1,
				name:  "",
				price: 1000,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewProduct(c.in.id, c.in.name, c.in.price)

			assert.Nil(t, got)
			assert.Equal(t, errors.Cause(err), domain.ErrInvalidInput)
		})
	}
}

func TestReCreateProduct(t *testing.T) {
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
			got := ReCreateProduct(c.in.id, c.in.name, c.in.price)
			assert.Equal(t, c.expected.id, got.ID())
			assert.Equal(t, c.expected.name, got.Name())
			assert.Equal(t, c.expected.price, got.Price())
		})
	}
}

func TestProduct_NewReview(t *testing.T) {
	type expected struct {
		id       ReviewID
		reviewTo ProductID
	}
	cases := []struct {
		name     string
		product  *Product
		in       ReviewID
		expected expected
	}{
		{
			name:    "正常系",
			product: ReCreateProduct(1, "product1", 100),
			in:      2,
			expected: expected{
				id:       2,
				reviewTo: 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.product.NewReview(c.in)
			assert.Equal(t, c.expected.id, got.ID())
			assert.Equal(t, c.expected.reviewTo, got.ReviewTo())
			assert.Zero(t, got.PostedBy())
			assert.Zero(t, got.Rating())
			assert.Zero(t, got.Comment())
		})
	}
}
