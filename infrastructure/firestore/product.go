package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/ryutah/virtual-ec/domain"
	"github.com/ryutah/virtual-ec/domain/model"
	"github.com/ryutah/virtual-ec/domain/repository"
)

var productErrMessages = struct {
	nextID          func(error) string
	get             func(model.ProductID, error) string
	getNoSuchEntity func(model.ProductID) string
	store           func(model.Product, error) string
}{
	nextID: func(err error) string {
		return fmt.Sprintf("failed to allocates id: %v", err)
	},
	get: func(id model.ProductID, err error) string {
		return fmt.Sprintf("failed to get product(%v): %v", id, err)
	},
	getNoSuchEntity: func(id model.ProductID) string {
		return fmt.Sprintf("product(%v) is not exists", id)
	},
	store: func(p model.Product, err error) string {
		return fmt.Sprintf("failed to store product(%v): %v", p, err)
	},
}

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
		return 0, errors.New(productErrMessages.nextID(err))
	}
	return model.ProductID(ids[0].ID), nil
}

func (p *Product) Get(ctx context.Context, id model.ProductID) (*model.Product, error) {
	var entity productEntity
	if err := p.client.Get(ctx, productKey(id), &entity); err == datastore.ErrNoSuchEntity {
		return nil, errors.Wrap(domain.ErrNoSuchEntity, productErrMessages.getNoSuchEntity(id))
	} else if err != nil {
		return nil, errors.New(productErrMessages.get(id, err))
	}

	return model.NewProduct(id, entity.Name, entity.Price), nil
}

func (p *Product) Store(ctx context.Context, product model.Product) error {
	entiity := productEntity{
		Name:  product.Name(),
		Price: product.Price(),
	}
	if _, err := p.client.Put(ctx, productKey(product.ID()), &entiity); err != nil {
		return errors.New(productErrMessages.store(product, err))
	}
	return nil
}

func productKey(id model.ProductID) *datastore.Key {
	return datastore.IDKey(kinds.product, int64(id), nil)
}

func (p *Product) Search(_ context.Context, _ repository.ProductQuery) (*repository.ProductSearchResult, error) {
	panic("not implemented") // TODO: Implement
}
