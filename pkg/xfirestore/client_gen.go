// Created by interfacer; DO NOT EDIT

package xfirestore

import (
	"context"

	"cloud.google.com/go/datastore"
)

// DatastoreClient is an interface generated for "cloud.google.com/go/datastore.Client".
type DatastoreClient interface {
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
	Run(context.Context, *datastore.Query) *datastore.Iterator
	RunInTransaction(context.Context, func(tx *datastore.Transaction) error, ...datastore.TransactionOption) (*datastore.Commit, error)
}
