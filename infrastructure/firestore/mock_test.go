package firestore_test

import (
	"context"
	"reflect"

	"cloud.google.com/go/datastore"
	. "github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/stretchr/testify/mock"
)

type clientStore interface {
	get(*datastore.Key) (v interface{}, ok bool)
	all() interface{}
}

type reviewSet struct {
	key *datastore.Key
	val *ReviewEntity
}

type reviewStore []reviewSet

var _ clientStore = (reviewStore)(nil)

func (r reviewStore) get(key *datastore.Key) (v interface{}, ok bool) {
	for _, set := range r {
		if set.key.Equal(key) {
			return set.val, true
		}
	}
	return nil, false
}

func (r reviewStore) all() interface{} {
	var entities []*ReviewEntity
	for _, set := range r {
		entities = append(entities, set.val)
	}
	return entities
}

type productSet struct {
	key *datastore.Key
	val *ProductEntity
}

type productStore []productSet

var _ productStore = (productStore)(nil)

func (p productStore) get(key *datastore.Key) (v interface{}, ok bool) {
	for _, set := range p {
		if set.key.Equal(key) {
			return set.val, true
		}
	}
	return nil, false
}

func (p productStore) all() interface{} {
	var entities []*ProductEntity
	for _, set := range p {
		entities = append(entities, set.val)
	}
	return entities
}

type mockClient struct {
	Client
	mock.Mock
	store clientStore
}

var _ Client = (*mockClient)(nil)

func (m *mockClient) withStore(store clientStore) *mockClient {
	m.store = store
	return m
}

func (m *mockClient) onAllocateIDs(ctx, keys interface{}) *mock.Call {
	return m.On("AllocateIDs", ctx, keys)
}

func (m *mockClient) onPut(ctx, key, v interface{}) *mock.Call {
	return m.On("Put", ctx, key, v)
}

func (m *mockClient) onGet(ctx, key, v interface{}) *mock.Call {
	return m.On("Get", ctx, key, v)
}

func (m *mockClient) onGetAll(ctx, query, v interface{}) *mock.Call {
	return m.On("GetAll", ctx, query, v)
}

func (m *mockClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	args := m.Called(ctx, keys)
	retKeys, ok := args.Get(0).([]*datastore.Key)
	if !ok {
		retKeys = nil
	}
	return retKeys, args.Error(1)
}

func (m *mockClient) Put(ctx context.Context, key *datastore.Key, v interface{}) (*datastore.Key, error) {
	args := m.Called(ctx, key, v)
	retKey, ok := args.Get(0).(*datastore.Key)
	if !ok {
		retKey = nil
	}
	return retKey, args.Error(1)
}

func (m *mockClient) Get(ctx context.Context, key *datastore.Key, v interface{}) error {
	args := m.Called(ctx, key, v)
	if m.store == nil {
		return args.Error(0)
	}

	got, ok := m.store.get(key)
	if !ok {
		return args.Error(0)
	}

	argVal, gotVal := reflect.ValueOf(v), reflect.ValueOf(got)
	argVal.Elem().Set(gotVal.Elem())
	return args.Error(0)
}

func (m *mockClient) GetAll(ctx context.Context, q *datastore.Query, v interface{}) ([]*datastore.Key, error) {
	args := m.Called(ctx, q, v)

	if m.store != nil {
		val := reflect.ValueOf(m.store.all())
		reflect.ValueOf(v).Elem().Set(val)
	}

	keys, _ := args.Get(0).([]*datastore.Key)
	return keys, args.Error(1)
}
