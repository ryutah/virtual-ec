package xfirestore_test

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	. "github.com/ryutah/virtual-ec/pkg/xfirestore"
	"github.com/stretchr/testify/mock"
)

type mockIterator struct {
	mock.Mock
	GenIterator
}

func (m *mockIterator) onCursor() *mock.Call {
	return m.On("Cursor")
}

func (m *mockIterator) onNext(v interface{}) *mock.Call {
	return m.On("Next", v)
}

func (m *mockIterator) Cursor() (datastore.Cursor, error) {
	args := m.Called()
	cursor, _ := args.Get(0).(datastore.Cursor)
	return cursor, args.Error(1)
}

func (m *mockIterator) Next(v interface{}) (*datastore.Key, error) {
	args := m.Called(v)
	key, _ := args.Get(0).(*datastore.Key)
	return key, args.Error(1)
}

type mockClient struct {
	DatastoreClient
	mock mock.Mock
}

func (m *mockClient) AssertExpectations(t *testing.T) bool {
	return m.mock.AssertExpectations(t)
}

func (m *mockClient) onAllocateIDs(ctx, datastoreKey interface{}) *mock.Call {
	return m.mock.On("AllocateIDs", ctx, datastoreKey)
}

func (m *mockClient) onClose() *mock.Call {
	return m.mock.On("Close")
}

func (m *mockClient) onCount(ctx, datastoreQuery interface{}) *mock.Call {
	return m.mock.On("Count", ctx, datastoreQuery)
}

func (m *mockClient) onDelete(ctx, datastoreKey interface{}) *mock.Call {
	return m.mock.On("Delete", ctx, datastoreKey)
}

func (m *mockClient) onDeleteMulti(ctx, datastoreKeys interface{}) *mock.Call {
	return m.mock.On("DeleteMulti", ctx, datastoreKeys)
}

func (m *mockClient) onGet(ctx, datastoreKey, v interface{}) *mock.Call {
	return m.mock.On("Get", ctx, datastoreKey, v)
}

func (m *mockClient) onGetAll(ctx, datastoreQuery, v interface{}) *mock.Call {
	return m.mock.On("GetAll", ctx, datastoreQuery, v)
}

func (m *mockClient) onGetMulti(ctx, datastoreKeys interface{}, v ...interface{}) *mock.Call {
	return m.mock.On("GetMulti", ctx, datastoreKeys, v)
}

func (m *mockClient) onMutate(ctx, datastoreMutations interface{}) *mock.Call {
	return m.mock.On("Mutate", ctx, datastoreMutations)
}

func (m *mockClient) onNewTransaction(ctx interface{}, datastoreTransactionOptions ...interface{}) *mock.Call {
	return m.mock.On("NewTransaction", ctx, datastoreTransactionOptions)
}

func (m *mockClient) onPut(ctx, datastoreKey, v interface{}) *mock.Call {
	return m.mock.On("Put", ctx, datastoreKey, v)
}

func (m *mockClient) onPutMulti(ctx, datastoreKeys, v interface{}) *mock.Call {
	return m.mock.On("PutMulti", ctx, datastoreKeys, v)
}

func (m *mockClient) onRun(ctx, datastoreQuery interface{}) *mock.Call {
	return m.mock.On("Run", ctx, datastoreQuery)
}

func (m *mockClient) onRunInTransaction(ctx, f, opts interface{}) *mock.Call {
	return m.mock.On("RunInTransaction", ctx, f, opts)
}

func (m *mockClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	args := m.mock.Called(ctx, keys)
	k, _ := args.Get(0).([]*datastore.Key)
	return k, args.Error(1)
}

func (m *mockClient) Close() error {
	return m.mock.Called().Error(0)
}

func (m *mockClient) Count(ctx context.Context, q *datastore.Query) (int, error) {
	args := m.mock.Called(ctx, q)
	return args.Int(0), args.Error(1)
}

func (m *mockClient) Delete(ctx context.Context, key *datastore.Key) error {
	return m.mock.Called(ctx, key).Error(0)
}

func (m *mockClient) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return m.mock.Called(ctx, keys).Error(0)
}

func (m *mockClient) Get(ctx context.Context, key *datastore.Key, v interface{}) error {
	return m.mock.Called(ctx, key, v).Error(0)
}

func (m *mockClient) GetAll(ctx context.Context, q *datastore.Query, v interface{}) ([]*datastore.Key, error) {
	args := m.mock.Called(ctx, q, v)
	keys, _ := args.Get(0).([]*datastore.Key)
	return keys, args.Error(1)
}

func (m *mockClient) GetMulti(ctx context.Context, keys []*datastore.Key, v interface{}) error {
	return m.mock.Called(ctx, keys, v).Error(0)
}

func (m *mockClient) Mutate(ctx context.Context, muts ...*datastore.Mutation) ([]*datastore.Key, error) {
	args := m.mock.Called(ctx, muts)
	keys, _ := args.Get(0).([]*datastore.Key)
	return keys, args.Error(1)
}

func (m *mockClient) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (*datastore.Transaction, error) {
	args := m.mock.Called(ctx, opts)
	tx, _ := args.Get(0).(*datastore.Transaction)
	return tx, args.Error(1)
}

func (m *mockClient) Put(ctx context.Context, key *datastore.Key, v interface{}) (*datastore.Key, error) {
	args := m.mock.Called(ctx, key, v)
	k, _ := args.Get(0).(*datastore.Key)
	return k, args.Error(1)
}

func (m *mockClient) PutMulti(ctx context.Context, keys []*datastore.Key, v interface{}) ([]*datastore.Key, error) {
	args := m.mock.Called(ctx, keys, v)
	ks, _ := args.Get(0).([]*datastore.Key)
	return ks, args.Error(1)
}

func (m *mockClient) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
	args := m.mock.Called(ctx, q)
	it, _ := args.Get(0).(*datastore.Iterator)
	return it
}

func (m *mockClient) RunInTransaction(ctx context.Context, f func(tx *datastore.Transaction) error, opts ...datastore.TransactionOption) (*datastore.Commit, error) {
	args := m.mock.Called(ctx, f, opts)
	c, _ := args.Get(0).(*datastore.Commit)
	return c, args.Error(1)
}
