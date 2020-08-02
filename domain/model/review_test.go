package model_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/domain/model"
	"gopkg.in/go-playground/assert.v1"
)

func TestNewReview(t *testing.T) {
	type (
		in struct {
			id       ReviewID
			reviewTo ProductID
			postedBy string
			rating   int
			comment  string
		}
		expected struct {
			id       ReviewID
			reviewTo ProductID
			postedBy string
			rating   int
			comment  string
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
				id:       1,
				reviewTo: 2,
				postedBy: "user1",
				rating:   3,
				comment:  "good!",
			},
			expected: expected{
				id:       1,
				reviewTo: 2,
				postedBy: "user1",
				rating:   3,
				comment:  "good!",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			review := NewReview(c.in.id, c.in.reviewTo, c.in.postedBy, c.in.rating, c.in.comment)
			assert.Equal(t, c.in.id, review.ID())
			assert.Equal(t, c.in.reviewTo, review.ReviewTo())
			assert.Equal(t, c.in.postedBy, review.PostedBy())
			assert.Equal(t, c.in.rating, review.Rating())
			assert.Equal(t, c.in.comment, review.Comment())
		})
	}
}
