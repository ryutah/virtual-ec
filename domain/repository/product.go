package repository

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
)

type Product interface {
	NextID(context.Context) (model.ProductID, error)
	Get(context.Context, model.ProductID) (*model.Product, error)
	Store(context.Context, model.Product) error
}
