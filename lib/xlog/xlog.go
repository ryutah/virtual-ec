package xlog

import (
	"context"
	"os"
	"sync"
)

type Logger interface {
	Debugf(ctx context.Context, format string, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})
	Warningf(ctx context.Context, format string, v ...interface{})
	Errorf(ctx context.Context, format string, v ...interface{})
}

var (
	logger Logger = newDefaultLogger(os.Stdout)
	mutex         = new(sync.Mutex)
)

func Register(l Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	logger = l
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(ctx, format, v...)
}

func Warningf(ctx context.Context, format string, v ...interface{}) {
	logger.Warningf(ctx, format, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(ctx, format, v...)
}
