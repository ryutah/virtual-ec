package firestore

var Kinds = struct {
	Product string
	Review  string
}{
	Product: kinds.product,
	Review:  kinds.review,
}

type (
	ProductEntity = productEntity
	ReviewEntity  = reviewEntity
)

var (
	ProductKey = productKey
	ReviewKey  = reviewKey
)
