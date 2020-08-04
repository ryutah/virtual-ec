package repository

import (
	"context"

	"github.com/ryutah/virtual-ec/domain/model"
)

type Review interface {
	NextID(context.Context) (model.ReviewID, error)
	Store(context.Context, model.Review) error
}
