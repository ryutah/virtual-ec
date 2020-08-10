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

func (p *Product) Get(ctx context.Context, id model.ProductID) (*model.Product, error) {
	var entity productEntity
	if err := p.client.Get(ctx, productKey(id), &entity); err != nil {
		return nil, errors.WithStack(err)
	}
	return model.NewProduct(id, entity.Name, entity.Price), nil
}

func (p *Product) Store(ctx context.Context, product model.Product) error {
	entiity := productEntity{
		Name:  product.Name(),
		Price: product.Price(),
	}
	_, err := p.client.Put(ctx, productKey(product.ID()), &entiity)
	return errors.WithStack(err)
}

func productKey(id model.ProductID) *datastore.Key {
	return datastore.IDKey(kinds.product, int64(id), nil)
}
