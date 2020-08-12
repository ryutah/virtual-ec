package xfirestore_test

import (
	"encoding/base64"
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	. "github.com/ryutah/virtual-ec/lib/xfirestore"
	"github.com/tj/assert"
)

func TestIterator_Cursor(t *testing.T) {
	genDummyCursor := func(s string) datastore.Cursor {
		b := base64.URLEncoding.EncodeToString([]byte(s))
		cursor, _ := datastore.DecodeCursor(b)
		return cursor
	}

	type (
		mocks struct {
			iterator_cursor_cursor datastore.Cursor
			iterator_cursor_err    error
		}
		expected struct {
			cursor datastore.Cursor
			err    error
		}
	)
	cases := []struct {
		name     string
		mocks    mocks
		expected expected
	}{
		{
			name: "正常なCursorを返す",
			mocks: mocks{
				iterator_cursor_cursor: genDummyCursor("dummy"),
				iterator_cursor_err:    nil,
			},
			expected: expected{
				cursor: genDummyCursor("dummy"),
				err:    nil,
			},
		},
		{
			name: "Cursor取得失敗",
			mocks: mocks{
				iterator_cursor_cursor: datastore.Cursor{},
				iterator_cursor_err:    errors.New("error"),
			},
			expected: expected{
				cursor: datastore.Cursor{},
				err:    errors.New("error"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inner := new(mockIterator)
			inner.onCursor().Return(
				c.mocks.iterator_cursor_cursor, c.mocks.iterator_cursor_err,
			)

			it := NewIterator(inner)
			got, err := it.Cursor()

			inner.AssertExpectations(t)
			assert.Equal(t, c.expected.cursor, got)
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func TestIterator_Next(t *testing.T) {
	type (
		mocks struct {
			iterator_cursor_next_datastoreKey *datastore.Key
			iterator_cursor_next_error        error
		}
		expected struct {
			args_iterator_cursor_next_v interface{}
			key                         *datastore.Key
			err                         error
		}
		dummy struct{}
	)
	cases := []struct {
		name     string
		in       *dummy
		mocks    mocks
		expected expected
	}{
		{
			name: "正常終了",
			in:   &dummy{},
			mocks: mocks{
				iterator_cursor_next_datastoreKey: datastore.IDKey("sample", 1, nil),
				iterator_cursor_next_error:        nil,
			},
			expected: expected{
				args_iterator_cursor_next_v: &dummy{},
				key:                         datastore.IDKey("sample", 1, nil),
				err:                         nil,
			},
		},
		{
			name: "エラー発生",
			in:   &dummy{},
			mocks: mocks{
				iterator_cursor_next_datastoreKey: nil,
				iterator_cursor_next_error:        errors.New("error"),
			},
			expected: expected{
				args_iterator_cursor_next_v: &dummy{},
				key:                         nil,
				err:                         errors.New("error"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inner := new(mockIterator)
			inner.onNext(c.expected.args_iterator_cursor_next_v).Return(
				c.mocks.iterator_cursor_next_datastoreKey, c.mocks.iterator_cursor_next_error,
			)

			it := NewIterator(inner)
			got, err := it.Next(c.in)

			inner.AssertExpectations(t)
			assert.Equal(t, c.expected.key, got)
			assert.Equal(t, c.expected.err, err)
		})
	}
}
