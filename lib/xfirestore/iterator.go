package xfirestore

import "cloud.google.com/go/datastore"

type Iterator interface {
	Cursor() (Cursor, error)
	Next(interface{}) (*datastore.Key, error)
}

type iteratorImpl struct {
	iterator iterator
}

func newIterator(it iterator) *iteratorImpl {
	return &iteratorImpl{
		iterator: it,
	}
}

var _ Iterator = (*iteratorImpl)(nil)

func (i *iteratorImpl) Cursor() (Cursor, error) {
	return i.iterator.Cursor()
}

func (i *iteratorImpl) Next(v interface{}) (*datastore.Key, error) {
	return i.iterator.Next(v)
}
