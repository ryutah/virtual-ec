package model

type ReviewID int

type Review struct {
	id       ReviewID
	reviewTo ProductID
	postedBy string
	rating   int
	comment  string
}

func NewReview(id ReviewID, reviewTo ProductID, postedBy string, rating int, comment string) *Review {
	// NOTE(ryutah): パラメータのバリデーションをすべき
	return &Review{
		id:       id,
		reviewTo: reviewTo,
		postedBy: postedBy,
		rating:   rating,
		comment:  comment,
	}
}

func (r *Review) ID() ReviewID {
	return r.id
}

func (r *Review) ReviewTo() ProductID {
	return r.reviewTo
}

func (r *Review) PostedBy() string {
	return r.postedBy
}

func (r *Review) Rating() int {
	return r.rating
}

func (r *Review) Comment() string {
	return r.comment
}
