// Created by interfacer; DO NOT EDIT

package internal

import (
	"context"

	"github.com/ryutah/virtual-ec/internal/usecase/admin"
)

// ProductCreator is an interface generated for "github.com/ryutah/virtual-ec/internal/usecase/admin.ProductCreate".
type ProductCreator interface {
	Create(context.Context, admin.ProductCreateInputPort, admin.ProductCreateOutputPort) bool
}
