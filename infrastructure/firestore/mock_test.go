package firestore_test

import (
	"context"
	"reflect"

	"cloud.google.com/go/datastore"
	. "github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/stretchr/testify/mock"
)

type mockClientStore map[datastore.Key]interface{}

type mockClient struct {
	Client
	mock.Mock
	store mockClientStore
}

var _ Client = (*mockClient)(nil)

func (m *mockClient) withStore(store mockClientStore) *mockClient {
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

	got, ok := m.store[*key]
	if !ok {
		return args.Error(0)
	}

	argVal, gotVal := reflect.ValueOf(v), reflect.ValueOf(got)
	argVal.Elem().Set(gotVal.Elem())
	return args.Error(0)
}
