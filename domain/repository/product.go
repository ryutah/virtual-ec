package repository

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
)

type Product interface {
	NextID(context.Context) (model.ProductID, error)
	Store(context.Context, model.Product) error
}
