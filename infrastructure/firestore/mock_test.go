package firestore_test

import (
	"context"
	"reflect"

	"cloud.google.com/go/datastore"
	. "github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/ryutah/virtual-ec/pkg/xfirestore"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"
)

type clientStore interface {
	get(*datastore.Key) (v interface{}, ok bool)
	all() interface{}
	allSet() []keyValueSet
}

type keyValueSet interface {
	key() *datastore.Key
	value() interface{}
}

type reviewSet struct {
	keyValueSet
	k *datastore.Key
	v *ReviewEntity
}

func (r reviewSet) key() *datastore.Key {
	return r.k
}

func (r reviewSet) value() interface{} {
	return r.v
}

type reviewStore []reviewSet

var _ clientStore = (reviewStore)(nil)

func (r reviewStore) get(key *datastore.Key) (v interface{}, ok bool) {
	for _, set := range r {
		if set.key().Equal(key) {
			return set.value(), true
		}
	}
	return nil, false
}

func (r reviewStore) all() interface{} {
	var entities []*ReviewEntity
	for _, set := range r {
		entities = append(entities, set.v)
	}
	return entities
}

func (r reviewStore) allSet() []keyValueSet {
	results := make([]keyValueSet, len(r))
	for i, v := range r {
		results[i] = v
	}
	return results
}

type productSet struct {
	k *datastore.Key
	v *ProductEntity
}

func (p productSet) key() *datastore.Key {
	return p.k
}

func (p productSet) value() interface{} {
	return p.v
}

type productStore []productSet

var _ productStore = (productStore)(nil)

func (p productStore) get(key *datastore.Key) (v interface{}, ok bool) {
	for _, set := range p {
		if set.key().Equal(key) {
			return set.value(), true
		}
	}
	return nil, false
}

func (p productStore) all() interface{} {
	var entities []*ProductEntity
	for _, set := range p {
		entities = append(entities, set.v)
	}
	return entities
}

func (p productStore) allSet() []keyValueSet {
	results := make([]keyValueSet, len(p))
	for i, v := range p {
		results[i] = v
	}
	return results
}

type mockClient struct {
	xfirestore.Client
	mock.Mock
	store clientStore
}

var _ xfirestore.Client = (*mockClient)(nil)

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

func (m *mockClient) onRun(ctx, datastoreQuery interface{}) *mock.Call {
	return m.On("Run", ctx, datastoreQuery)
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

func (m *mockClient) Run(ctx context.Context, q *datastore.Query) xfirestore.Iterator {
	args := m.Called(ctx, q)
	it, _ := args.Get(0).(xfirestore.Iterator)
	return it
}

type mockIterator struct {
	mock.Mock
	xfirestore.Iterator
	store  []keyValueSet
	cursor int
}

func (m *mockIterator) withStore(store clientStore) *mockIterator {
	m.store = store.allSet()
	return m
}

// NOTE(ryutah): ReturnはデフォルトでReturn(nil, nil)が設定される
// エラーを返すようにしたい場合呼び出し元でReturn(nil, err)と設定すること(第一引数は設定しても無視される)
func (m *mockIterator) onNext(v interface{}) *mock.Call {
	return m.On("Next", v).Return(nil, nil)
}

func (m *mockIterator) Next(v interface{}) (*datastore.Key, error) {
	args := m.Called(v)
	if err := args.Error(1); err != nil {
		return nil, err
	}

	if m.cursor >= len(m.store) {
		return nil, iterator.Done
	}
	setValue(v, m.store[m.cursor].value())
	key := m.store[m.cursor].key()
	m.cursor++
	return key, args.Error(1)
}

func setValue(target, src interface{}) {
	tgtVal, srcVal := reflect.ValueOf(target), reflect.ValueOf(src)
	tgtVal.Elem().Set(srcVal.Elem())
}
