package xlog

import (
	"context"
	"fmt"
	"io"
	"log"
)

type defaultLogger struct {
	logger *log.Logger
}

var _ Logger = (*defaultLogger)(nil)

func newDefaultLogger(w io.Writer) *defaultLogger {
	return &defaultLogger{
		logger: log.New(w, "", 0),
	}
}

func (l *defaultLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	format = fmt.Sprintf("[DEBUG] %s", format)
	l.logger.Printf(format, v...)
}

func (l *defaultLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	format = fmt.Sprintf("[INFO] %s", format)
	l.logger.Printf(format, v...)
}

func (l *defaultLogger) Warningf(ctx context.Context, format string, v ...interface{}) {
	format = fmt.Sprintf("[WARN] %s", format)
	l.logger.Printf(format, v...)
}

func (l *defaultLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	format = fmt.Sprintf("[ERROR] %s", format)
	l.logger.Printf(format, v...)
}
