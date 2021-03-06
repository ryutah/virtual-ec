// Created by interfacer; DO NOT EDIT

package xfirestore

import (
	"cloud.google.com/go/datastore"
)

// transaction is an interface generated for "cloud.google.com/go/datastore.Transaction".
type transaction interface {
	Commit() (*datastore.Commit, error)
	Delete(*datastore.Key) error
	DeleteMulti([]*datastore.Key) error
	Get(*datastore.Key, interface{}) error
	GetMulti([]*datastore.Key, interface{}) error
	Mutate(...*datastore.Mutation) ([]*datastore.PendingKey, error)
	Put(*datastore.Key, interface{}) (*datastore.PendingKey, error)
	PutMulti([]*datastore.Key, interface{}) ([]*datastore.PendingKey, error)
	Rollback() error
}
