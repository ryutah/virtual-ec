package xlog_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/ryutah/virtual-ec/lib/xlog"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	type in struct {
		format string
		v      []interface{}
	}
	cases := []struct {
		name     string
		in       in
		method   func(*DefaultLogger, context.Context, string, ...interface{})
		expected string
	}{
		{
			name: "Debugf",
			in: in{
				format: "this is debug %v",
				v:      []interface{}{"value"},
			},
			method:   (*DefaultLogger).Debugf,
			expected: "[DEBUG] this is debug value\n",
		},
		{
			name: "Infof",
			in: in{
				format: "this is info %v",
				v:      []interface{}{"value"},
			},
			method:   (*DefaultLogger).Infof,
			expected: "[INFO] this is info value\n",
		},
		{
			name: "Warningf",
			in: in{
				format: "this is warning %v",
				v:      []interface{}{"value"},
			},
			method:   (*DefaultLogger).Warningf,
			expected: "[WARN] this is warning value\n",
		},
		{
			name: "Errorf",
			in: in{
				format: "this is error %v",
				v:      []interface{}{"value"},
			},
			method:   (*DefaultLogger).Errorf,
			expected: "[ERROR] this is error value\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := NewDefualtLogger(buf)
			c.method(logger, context.Background(), c.in.format, c.in.v...)
			assert.Equal(t, c.expected, buf.String())
		})
	}
}
