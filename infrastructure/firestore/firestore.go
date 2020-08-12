//go:generate interfacer -for cloud.google.com/go/datastore.Client -as firestore.Client -o client.go
//go:generate interfacer -for cloud.google.com/go/datastore.Iterator -as firestore.Iterator -o iterator.go
//go:generate interfacer -for cloud.google.com/go/datastore.Cursor -as firestore.Cursor -o cursor.go
//go:generate interfacer -for cloud.google.com/go/datastore.Transaction -as firestore.Transaction -o transaction.go

package firestore

import (
	"context"

	"cloud.google.com/go/datastore"
)

var kinds = struct {
	product string
	review  string
}{
	product: "product",
	review:  "review",
}

type ClientWrapper interface {
	AllocateIDs(context.Context, []*datastore.Key) ([]*datastore.Key, error)
	Close() error
	Count(context.Context, *datastore.Query) (int, error)
	Delete(context.Context, *datastore.Key) error
	DeleteMulti(context.Context, []*datastore.Key) error
	Get(context.Context, *datastore.Key, interface{}) error
	GetAll(context.Context, *datastore.Query, interface{}) ([]*datastore.Key, error)
	GetMulti(context.Context, []*datastore.Key, interface{}) error
	Mutate(context.Context, ...*datastore.Mutation) ([]*datastore.Key, error)
	NewTransaction(context.Context, ...datastore.TransactionOption) (*datastore.Transaction, error)
	Put(context.Context, *datastore.Key, interface{}) (*datastore.Key, error)
	PutMulti(context.Context, []*datastore.Key, interface{}) ([]*datastore.Key, error)
	Run(context.Context, *datastore.Query) IteratorWrapper
	RunInTransaction(context.Context, func(tx *datastore.Transaction) error, ...datastore.TransactionOption) (*datastore.Commit, error)
}

type clientWrapperImpl struct {
	client Client
}

func (c *clientWrapperImpl) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	return c.client.AllocateIDs(ctx, keys)
}

func (c *clientWrapperImpl) Close() error {
	return c.client.Close()
}

func (c *clientWrapperImpl) Count(ctx context.Context, q *datastore.Query) (int, error) {
	return c.client.Count(ctx, q)
}

func (c *clientWrapperImpl) Delete(ctx context.Context, key *datastore.Key) error {
	return c.client.Delete(ctx, key)
}

func (c *clientWrapperImpl) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return c.client.DeleteMulti(ctx, keys)
}

func (c *clientWrapperImpl) Get(ctx context.Context, key *datastore.Key, v interface{}) error {
	return c.client.Get(ctx, key, v)
}

func (c *clientWrapperImpl) GetAll(ctx context.Context, q *datastore.Query, v interface{}) ([]*datastore.Key, error) {
	return c.client.GetAll(ctx, q, v)
}

func (c *clientWrapperImpl) GetMulti(ctx context.Context, keys []*datastore.Key, v interface{}) error {
	return c.client.GetMulti(ctx, keys, v)
}

func (c *clientWrapperImpl) Mutate(ctx context.Context, m ...*datastore.Mutation) ([]*datastore.Key, error) {
	return c.client.Mutate(ctx, m...)
}

func (c *clientWrapperImpl) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (*datastore.Transaction, error) {
	return c.client.NewTransaction(ctx, opts...)
}

func (c *clientWrapperImpl) Put(ctx context.Context, key *datastore.Key, v interface{}) (*datastore.Key, error) {
	return c.client.Put(ctx, key, v)
}

func (c *clientWrapperImpl) PutMulti(ctx context.Context, keys []*datastore.Key, v interface{}) ([]*datastore.Key, error) {
	return c.client.PutMulti(ctx, keys, v)
}

func (c *clientWrapperImpl) Run(ctx context.Context, q *datastore.Query) IteratorWrapper {
	it := c.client.Run(ctx, q)
	return &iteratorWrapperImpl{iterator: it}
}

func (c *clientWrapperImpl) RunInTransaction(ctx context.Context, f func(tx *datastore.Transaction) error, opts ...datastore.TransactionOption) (*datastore.Commit, error) {
	return c.client.RunInTransaction(ctx, f, opts...)
}

type IteratorWrapper interface {
	Cursor() (Cursor, error)
	Next(interface{}) (*datastore.Key, error)
}

type iteratorWrapperImpl struct {
	iterator Iterator
}

func (i *iteratorWrapperImpl) Cursor() (Cursor, error) {
	return i.iterator.Cursor()
}

func (i *iteratorWrapperImpl) Next(v interface{}) (*datastore.Key, error) {
	return i.iterator.Next(v)
}
