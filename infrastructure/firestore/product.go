package firestore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

type productEntity struct {
	Name  string
	Price int
}

type Product struct {
	client Client
}

var _ repository.Product = (*Product)(nil)

func NewProduct(client Client) *Product {
	return &Product{
		client: client,
	}
}

func (p *Product) NextID(ctx context.Context) (model.ProductID, error) {
	ids, err := p.client.AllocateIDs(ctx, []*datastore.Key{
		datastore.IncompleteKey(kinds.product, nil),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return model.ProductID(ids[0].ID), nil
}

func (p *Product) Store(ctx context.Context, product model.Product) error {
	entiity := productEntity{
		Name:  product.Name(),
		Price: product.Price(),
	}
	key := datastore.IDKey(kinds.product, int64(product.ID()), nil)
	_, err := p.client.Put(ctx, key, &entiity)
	return errors.WithStack(err)
}
