package xfirestore_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	. "github.com/ryutah/virtual-ec/pkg/xfirestore"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"cloud.google.com/go/datastore"
)

func TestClient_WrapperedFunctions(t *testing.T) {
	type dummy struct {
		Name string
	}
	dummySlicePtr := func(s []dummy) *[]dummy { return &s }

	ctx := context.Background()

	cases := []struct {
		name     string
		funcName string
		in       []interface{}
		args     []interface{} // Optional. 可変長引数受け取る関数用
		expected []interface{}
	}{
		{
			name:     "AllocateIDs",
			funcName: "AllocateIDs",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IncompleteKey("sample", nil),
					datastore.IncompleteKey("sample2", nil),
					datastore.IncompleteKey("sample3", datastore.IDKey("sample4", 1, nil)),
				},
			},
			expected: []interface{}{
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, datastore.IDKey("sample4", 1, nil)),
				},
				nil,
			},
		},
		{
			name:     "AllocateIDs_Error",
			funcName: "AllocateIDs",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IncompleteKey("sample", nil),
					datastore.IncompleteKey("sample2", nil),
					datastore.IncompleteKey("sample3", datastore.IDKey("sample4", 1, nil)),
				},
			},
			expected: []interface{}{
				([]*datastore.Key)(nil), errors.New("error"),
			},
		},
		{
			name:     "Close",
			funcName: "Close",
			in:       []interface{}{},
			expected: []interface{}{
				nil,
			},
		},
		{
			name:     "Close_Error",
			funcName: "Close",
			in:       []interface{}{},
			expected: []interface{}{
				errors.New("error"),
			},
		},
		{
			name:     "Count",
			funcName: "Count",
			in: []interface{}{
				ctx, datastore.NewQuery("sample"),
			},
			expected: []interface{}{
				10, nil,
			},
		},
		{
			name:     "Count_Error",
			funcName: "Count",
			in: []interface{}{
				ctx, datastore.NewQuery("sample"),
			},
			expected: []interface{}{
				0, errors.New("error"),
			},
		},
		{
			name:     "Delete",
			funcName: "Delete",
			in: []interface{}{
				ctx, datastore.IDKey("sample", 1, nil),
			},
			expected: []interface{}{
				nil,
			},
		},
		{
			name:     "Delete_Error",
			funcName: "Delete",
			in: []interface{}{
				ctx, datastore.IDKey("sample", 1, nil),
			},
			expected: []interface{}{
				errors.New("error"),
			},
		},
		{
			name:     "DeleteMulti",
			funcName: "DeleteMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
			},
			expected: []interface{}{
				nil,
			},
		},
		{
			name:     "DeleteMulti_Error",
			funcName: "DeleteMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
			},
			expected: []interface{}{
				errors.New("error"),
			},
		},
		{
			name:     "Get",
			funcName: "Get",
			in: []interface{}{
				ctx,
				datastore.IDKey("sample", 1, nil),
				&dummy{},
			},
			expected: []interface{}{
				nil,
			},
		},
		{
			name:     "Get_Error",
			funcName: "Get",
			in: []interface{}{
				ctx,
				datastore.IDKey("sample", 1, nil),
				&dummy{},
			},
			expected: []interface{}{
				errors.New("error"),
			},
		},
		{
			name:     "GetAll",
			funcName: "GetAll",
			in: []interface{}{
				ctx,
				datastore.NewQuery("sample"),
				dummySlicePtr(make([]dummy, 0)),
			},
			expected: []interface{}{
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				nil,
			},
		},
		{
			name:     "GetAll_Error",
			funcName: "GetAll",
			in: []interface{}{
				ctx,
				datastore.NewQuery("sample"),
				dummySlicePtr(make([]dummy, 0)),
			},
			expected: []interface{}{
				([]*datastore.Key)(nil),
				errors.New("error"),
			},
		},
		{
			name:     "GetMulti",
			funcName: "GetMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				make([]dummy, 3),
			},
			expected: []interface{}{
				nil,
			},
		},
		{
			name:     "GetMulti_Error",
			funcName: "GetMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				make([]dummy, 3),
			},
			expected: []interface{}{
				errors.New("Error"),
			},
		},
		{
			name:     "Mutate",
			funcName: "Mutate",
			in: []interface{}{
				ctx,
				datastore.NewInsert(datastore.IDKey("sample", 1, nil), &dummy{Name: "name"}),
				datastore.NewUpdate(datastore.IDKey("sample2", 2, nil), &dummy{Name: "name2"}),
				datastore.NewDelete(datastore.IDKey("sample3", 3, nil)),
			},
			args: []interface{}{
				ctx,
				[]*datastore.Mutation{
					datastore.NewInsert(datastore.IDKey("sample", 1, nil), &dummy{Name: "name"}),
					datastore.NewUpdate(datastore.IDKey("sample2", 2, nil), &dummy{Name: "name2"}),
					datastore.NewDelete(datastore.IDKey("sample3", 3, nil)),
				},
			},
			expected: []interface{}{
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				nil,
			},
		},
		{
			name:     "Mutate_Error",
			funcName: "Mutate",
			in: []interface{}{
				ctx,
				datastore.NewInsert(datastore.IDKey("sample", 1, nil), &dummy{Name: "name"}),
				datastore.NewUpdate(datastore.IDKey("sample2", 2, nil), &dummy{Name: "name2"}),
				datastore.NewDelete(datastore.IDKey("sample3", 3, nil)),
			},
			args: []interface{}{
				ctx,
				[]*datastore.Mutation{
					datastore.NewInsert(datastore.IDKey("sample", 1, nil), &dummy{Name: "name"}),
					datastore.NewUpdate(datastore.IDKey("sample2", 2, nil), &dummy{Name: "name2"}),
					datastore.NewDelete(datastore.IDKey("sample3", 3, nil)),
				},
			},
			expected: []interface{}{
				([]*datastore.Key)(nil), errors.New("error"),
			},
		},
		{
			name:     "Put",
			funcName: "Put",
			in: []interface{}{
				ctx,
				datastore.IncompleteKey("sample", nil),
				&dummy{Name: "name"},
			},
			expected: []interface{}{
				datastore.IDKey("sample", 1, nil), nil,
			},
		},
		{
			name:     "Put_Error",
			funcName: "Put",
			in: []interface{}{
				ctx,
				datastore.IncompleteKey("sample", nil),
				&dummy{Name: "name"},
			},
			expected: []interface{}{
				(*datastore.Key)(nil), errors.New("error"),
			},
		},
		{
			name:     "PutMulti",
			funcName: "PutMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IncompleteKey("sample", nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				[]*dummy{
					{Name: "name"},
					{Name: "name2"},
					{Name: "name3"},
				},
			},
			expected: []interface{}{
				[]*datastore.Key{
					datastore.IDKey("sample", 1, nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				nil,
			},
		},
		{
			name:     "PutMulti_Error",
			funcName: "PutMulti",
			in: []interface{}{
				ctx,
				[]*datastore.Key{
					datastore.IncompleteKey("sample", nil),
					datastore.IDKey("sample2", 2, nil),
					datastore.IDKey("sample3", 3, nil),
				},
				[]*dummy{
					{Name: "name"},
					{Name: "name2"},
					{Name: "name3"},
				},
			},
			expected: []interface{}{
				([]*datastore.Key)(nil), errors.New("error"),
			},
		},
		{
			name:     "NewTransaction",
			funcName: "NewTransaction",
			in: []interface{}{
				ctx,
				datastore.MaxAttempts(1),
				datastore.ReadOnly,
			},
			args: []interface{}{
				ctx,
				[]datastore.TransactionOption{
					datastore.MaxAttempts(1),
					datastore.ReadOnly,
				},
			},
			expected: []interface{}{
				new(datastore.Transaction), nil,
			},
		},
		{
			name:     "NewTransaction_Error",
			funcName: "NewTransaction",
			in: []interface{}{
				ctx,
				datastore.MaxAttempts(1),
				datastore.ReadOnly,
			},
			args: []interface{}{
				ctx,
				[]datastore.TransactionOption{
					datastore.MaxAttempts(1),
					datastore.ReadOnly,
				},
			},
			expected: []interface{}{
				(*datastore.Transaction)(nil), errors.New("error"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inner := new(mockClient)
			if c.args == nil {
				inner.mock.On(c.funcName, c.in...).Return(c.expected...)
			} else {
				inner.mock.On(c.funcName, c.args...).Return(c.expected...)
			}

			client := NewClient(inner)
			var args []reflect.Value
			for _, arg := range c.in {
				args = append(args, reflect.ValueOf(arg))
			}
			rets := reflect.ValueOf(client).MethodByName(c.funcName).Call(args)

			inner.AssertExpectations(t)
			for i, ret := range rets {
				assert.Equal(t, c.expected[i], ret.Interface())
			}
		})
	}
}

func TestClient_Run(t *testing.T) {
	ctx := context.Background()

	inner := new(mockClient)
	inner.onRun(ctx, datastore.NewQuery("sample")).Return(new(datastore.Iterator))

	client := NewClient(inner)
	got := client.Run(ctx, datastore.NewQuery("sample"))

	inner.AssertExpectations(t)
	assert.Equal(t, NewIterator(new(datastore.Iterator)), got)
}

func TestClient_RunInTransaction(t *testing.T) {
	ctx := context.Background()

	var called bool
	innerFunc := func(tx Transaction) error {
		assert.Equal(t, tx, new(datastore.Transaction))
		called = true
		return nil
	}

	inner := new(mockClient)
	inner.
		onRunInTransaction(ctx, mock.Anything, []datastore.TransactionOption{
			datastore.MaxAttempts(1), datastore.ReadOnly},
		).
		Return(new(datastore.Commit), errors.New("error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(*datastore.Transaction) error)
			f(new(datastore.Transaction))
		})

	client := NewClient(inner)
	got, err := client.RunInTransaction(ctx, innerFunc, datastore.MaxAttempts(1), datastore.ReadOnly)

	inner.AssertExpectations(t)
	assert.True(t, called)
	assert.Equal(t, new(datastore.Commit), got)
	assert.Equal(t, errors.New("error"), err)
}
