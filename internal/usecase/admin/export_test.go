package admin

import "github.com/ryutah/virtual-ec/internal/domain/model"

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

var ProductSearchErrorMessages = struct {
	Failed func() string
}{
	Failed: productSearchErrorMessages.failed,
}
