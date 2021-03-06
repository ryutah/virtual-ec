// Created by interfacer; DO NOT EDIT

package internal

import (
	"context"

	"github.com/ryutah/virtual-ec/internal/usecase/admin"
)

// ProductSearcher is an interface generated for "github.com/ryutah/virtual-ec/internal/usecase/admin.ProductSearch".
type ProductSearcher interface {
	Search(context.Context, admin.ProductSearchInputPort, admin.ProductSearchOutputPort) bool
}
