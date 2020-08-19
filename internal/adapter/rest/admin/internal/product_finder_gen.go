// Created by interfacer; DO NOT EDIT

package internal

import (
	"context"

	"github.com/ryutah/virtual-ec/internal/usecase/admin"
)

// ProductFinder is an interface generated for "github.com/ryutah/virtual-ec/internal/usecase/admin.ProductFind".
type ProductFinder interface {
	Find(context.Context, int, admin.ProductFindOutputPort) bool
}
