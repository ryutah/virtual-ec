package repository_test

import (
	"testing"

	"github.com/ryutah/virtual-ec/domain/model"
	. "github.com/ryutah/virtual-ec/domain/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewReviewQuery(t *testing.T) {
	type expected struct {
		productID    model.ProductID
		productID_ok bool
	}
	cases := []struct {
		name     string
		callers  []func(ReviewQuery) ReviewQuery
		expected expected
	}{
		{
			name:    "生成のみ",
			callers: nil,
			expected: expected{
				productID:    0,
				productID_ok: false,
			},
		},
		{
			name: "プロダクトID指定",
			callers: []func(ReviewQuery) ReviewQuery{
				func(q ReviewQuery) ReviewQuery { return q.WithProductID(1) },
			},
			expected: expected{
				productID:    1,
				productID_ok: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query := NewReviewQuery()
			for _, f := range c.callers {
				query = f(query)
			}
			gotProductID, ok := query.ProductID()
			assert.Equal(t, c.expected.productID, gotProductID)
			assert.Equal(t, c.expected.productID_ok, ok)
		})
	}
}
