package admin

import "github.com/ryutah/virtual-ec/domain/model"

var ProductFindFailedErrorMessages = struct {
	NotFound func(model.ProductID) string
	Failed   func(model.ProductID) string
}{
	NotFound: productFindFailedErrorMessages.notFound,
	Failed:   productFindFailedErrorMessages.failed,
}

var ProductCreateErrroMessages = struct {
	Failed func() string
}{
	Failed: productCreateErrroMessages.failed,
}
