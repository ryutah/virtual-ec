package xfirestore

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Client interface {
	AllocateIDs(context.Context, []*datastore.Key) ([]*datastore.Key, error)
	Close() error
	Count(context.Context, *datastore.Query) (int, error)
	Delete(context.Context, *datastore.Key) error
	DeleteMulti(context.Context, []*datastore.Key) error
	Get(context.Context, *datastore.Key, interface{}) error
	GetAll(context.Context, *datastore.Query, interface{}) ([]*datastore.Key, error)
	GetMulti(context.Context, []*datastore.Key, interface{}) error
	Mutate(context.Context, ...*datastore.Mutation) ([]*datastore.Key, error)
	NewTransaction(context.Context, ...datastore.TransactionOption) (Transaction, error)
	Put(context.Context, *datastore.Key, interface{}) (*datastore.Key, error)
	PutMulti(context.Context, []*datastore.Key, interface{}) ([]*datastore.Key, error)
	Run(context.Context, *datastore.Query) Iterator
	RunInTransaction(context.Context, func(tx Transaction) error, ...datastore.TransactionOption) (*datastore.Commit, error)
}

type clientImpl struct {
	client DatastoreClient
}

func NewClient(c DatastoreClient) Client {
	return &clientImpl{
		client: c,
	}
}

func (c *clientImpl) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	return c.client.AllocateIDs(ctx, keys)
}

func (c *clientImpl) Close() error {
	return c.client.Close()
}

func (c *clientImpl) Count(ctx context.Context, q *datastore.Query) (int, error) {
	return c.client.Count(ctx, q)
}

func (c *clientImpl) Delete(ctx context.Context, key *datastore.Key) error {
	return c.client.Delete(ctx, key)
}

func (c *clientImpl) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return c.client.DeleteMulti(ctx, keys)
}

func (c *clientImpl) Get(ctx context.Context, key *datastore.Key, v interface{}) error {
	return c.client.Get(ctx, key, v)
}

func (c *clientImpl) GetAll(ctx context.Context, q *datastore.Query, v interface{}) ([]*datastore.Key, error) {
	return c.client.GetAll(ctx, q, v)
}

func (c *clientImpl) GetMulti(ctx context.Context, keys []*datastore.Key, v interface{}) error {
	return c.client.GetMulti(ctx, keys, v)
}

func (c *clientImpl) Mutate(ctx context.Context, muts ...*datastore.Mutation) ([]*datastore.Key, error) {
	return c.client.Mutate(ctx, muts...)
}

func (c *clientImpl) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (Transaction, error) {
	return c.client.NewTransaction(ctx, opts...)
}

func (c *clientImpl) Put(ctx context.Context, key *datastore.Key, v interface{}) (*datastore.Key, error) {
	return c.client.Put(ctx, key, v)
}

func (c *clientImpl) PutMulti(ctx context.Context, keys []*datastore.Key, v interface{}) ([]*datastore.Key, error) {
	return c.client.PutMulti(ctx, keys, v)
}

func (c *clientImpl) Run(ctx context.Context, q *datastore.Query) Iterator {
	it := c.client.Run(ctx, q)
	return newIterator(it)
}

func (c *clientImpl) RunInTransaction(ctx context.Context, f func(tx Transaction) error, opts ...datastore.TransactionOption) (*datastore.Commit, error) {
	return c.client.RunInTransaction(ctx, wrapRunInTransactionFunc(f), opts...)
}

func wrapRunInTransactionFunc(f func(Transaction) error) func(*datastore.Transaction) error {
	return func(tx *datastore.Transaction) error {
		return f(tx)
	}
}
