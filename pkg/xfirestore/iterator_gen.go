// Created by interfacer; DO NOT EDIT

package xfirestore

import (
	"cloud.google.com/go/datastore"
)

// iterator is an interface generated for "cloud.google.com/go/datastore.Iterator".
type iterator interface {
	Cursor() (datastore.Cursor, error)
	Next(interface{}) (*datastore.Key, error)
}
