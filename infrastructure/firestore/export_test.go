package firestore

import (
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

var Kinds = struct {
	Product string
	Review  string
}{
	Product: kinds.product,
	Review:  kinds.review,
}

var (
	ProductErrMessages = struct {
		NextID          func(error) string
		Get             func(model.ProductID, error) string
		GetNoSuchEntity func(model.ProductID) string
		Store           func(model.Product, error) string
		Search          func(err error) string
	}{
		NextID:          productErrMessages.nextID,
		Get:             productErrMessages.get,
		GetNoSuchEntity: productErrMessages.getNoSuchEntity,
		Store:           productErrMessages.store,
		Search:          productErrMessages.search,
	}

	ReviewErrMessages = struct {
		NextID func(error) string
		Store  func(model.Review, error) string
		Search func(repository.ReviewQuery, error) string
	}{
		NextID: reviewErrMessages.nextID,
		Store:  reviewErrMessages.store,
		Search: reviewErrMessages.search,
	}
)

type (
	ProductEntity = productEntity
	ReviewEntity  = reviewEntity
)

var (
	ProductKey = productKey
	ReviewKey  = reviewKey
)
