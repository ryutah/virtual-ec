package consumer

import "github.com/ryutah/virtual-ec/domain/model"

var ReviewListErrorMessages = struct {
	Failed func(model.ProductID) string
}{
	Failed: reviewListErrorMessages.failed,
}

var ReviewPostErrorMessages = struct {
	ProductNotFound func(model.ProductID) string
	Failed          func() string
}{
	ProductNotFound: reviewPostErrorMessages.productNotFound,
	Failed:          reviewPostErrorMessages.failed,
}
