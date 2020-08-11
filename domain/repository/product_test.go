package repository_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/domain/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductQuery(t *testing.T) {
	type (
		caller   func(ProductQuery) ProductQuery
		expected struct {
			name    string
			name_ok bool
		}
	)
	cases := []struct {
		name     string
		callers  []caller
		expected expected
	}{
		{
			name:     "検索条件未指定",
			callers:  []caller{},
			expected: expected{},
		},
		{
			name: "Nameを指定",
			callers: []caller{
				func(p ProductQuery) ProductQuery {
					return p.WithName("product")
				},
			},
			expected: expected{
				name:    "product",
				name_ok: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query := NewProductQuery()
			for _, ca := range c.callers {
				query = ca(query)
			}

			gotName, nameOK := query.Name()
			assert.Equal(t, c.expected.name, gotName)
			assert.Equal(t, c.expected.name_ok, nameOK)
		})
	}
}
