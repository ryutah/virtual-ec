package xlog_test

import (
	"context"
	"testing"

	. "github.com/ryutah/virtual-ec/lib/xlog"
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

var _ Logger = (*mockLogger)(nil)

func (m *mockLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	m.Called(ctx, format, v)
}

func (m *mockLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	m.Called(ctx, format, v)
}

func (m *mockLogger) Warningf(ctx context.Context, format string, v ...interface{}) {
	m.Called(ctx, format, v)
}

func (m *mockLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	m.Called(ctx, format, v)
}

func TestDebugf(t *testing.T) {
	ctx := context.Background()

	logger := new(mockLogger)
	logger.On("Debugf", ctx, "log %v", []interface{}{"value"})

	Register(logger)
	Debugf(ctx, "log %v", "value")

	logger.AssertExpectations(t)
}

func TestInfof(t *testing.T) {
	ctx := context.Background()

	logger := new(mockLogger)
	logger.On("Infof", ctx, "log %v", []interface{}{"value"})

	Register(logger)
	Infof(ctx, "log %v", "value")

	logger.AssertExpectations(t)
}

func TestWarningf(t *testing.T) {
	ctx := context.Background()

	logger := new(mockLogger)
	logger.On("Warningf", ctx, "log %v", []interface{}{"value"})

	Register(logger)
	Warningf(ctx, "log %v", "value")

	logger.AssertExpectations(t)
}

func TestErrorf(t *testing.T) {
	ctx := context.Background()

	logger := new(mockLogger)
	logger.On("Errorf", ctx, "log %v", []interface{}{"value"})

	Register(logger)
	Errorf(ctx, "log %v", "value")

	logger.AssertExpectations(t)
}
