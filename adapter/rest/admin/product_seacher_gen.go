// Code generated by interfacer; DO NOT EDIT

package admin

import (
	"context"

	"github.com/ryutah/virtual-ec/usecase/admin"
)

// ProductSearcher is an interface generated for "github.com/ryutah/virtual-ec/usecase/admin.ProductSearch".
type ProductSearcher interface {
	Search(context.Context, admin.ProductSearchInputPort, admin.ProductSearchOutputPort) bool
}
