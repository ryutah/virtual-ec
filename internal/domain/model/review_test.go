package model_test

import (
	"testing"

	"github.com/ryutah/virtual-ec/internal/domain/model"
	. "github.com/ryutah/virtual-ec/internal/domain/model"
	"gopkg.in/go-playground/assert.v1"
)

func TestReCreateReview(t *testing.T) {
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
			review := ReCreateReview(c.in.id, c.in.reviewTo, c.in.postedBy, c.in.rating, c.in.comment)
			assert.Equal(t, c.in.id, review.ID())
			assert.Equal(t, c.in.reviewTo, review.ReviewTo())
			assert.Equal(t, c.in.postedBy, review.PostedBy())
			assert.Equal(t, c.in.rating, review.Rating())
			assert.Equal(t, c.in.comment, review.Comment())
		})
	}
}

func TestReview_Write(t *testing.T) {
	type (
		in struct {
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
		review   *Review
		in       in
		expected expected
	}{
		{
			name:   "正常系",
			review: model.ReCreateProduct(1, "p1", 100).NewReview(2),
			in: in{
				postedBy: "user1",
				rating:   3,
				comment:  "Good!",
			},
			expected: expected{
				id:       2,
				reviewTo: 1,
				postedBy: "user1",
				rating:   3,
				comment:  "Good!",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.review.Write(c.in.postedBy, c.in.rating, c.in.comment)
			assert.Equal(t, c.expected.id, c.review.ID())
			assert.Equal(t, c.expected.reviewTo, c.review.ReviewTo())
			assert.Equal(t, c.expected.postedBy, c.review.PostedBy())
			assert.Equal(t, c.expected.comment, c.review.Comment())
			assert.Equal(t, c.expected.rating, c.review.Rating())
		})
	}
}
