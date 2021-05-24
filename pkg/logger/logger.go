package logger

import (
	"fmt"
	"io"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (

	// Logger to print stuff
	Logger interface {
		Info(msg string, keysAndValues ...interface{})
		V(level int) Logger
		Fork(keysAndValues ...interface{}) Logger
		DF() io.WriteCloser
	}

	// Options to create new Logger
	Options struct {
		WriterHandlers []io.Writer
		WithCaller     bool
		WithStacktrace bool
		JSONFormat     bool
	}

	logger struct {
		lgr logr.Logger
		fd  io.WriteCloser
	}
)

// New builds Logger
func New(options Options) Logger {
	var zlog *zap.Logger
	if options.JSONFormat {
		z, err := zap.NewProduction(zap.AddCallerSkip(1))
		if err != nil {
			panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
		}
		zlog = z
	} else {
		z, err := zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
		}
		zlog = z

	}
	if options.WithCaller {
		zlog.WithOptions(zap.WithCaller(true))
	}
	if options.WithStacktrace {
		zlog.WithOptions(zap.AddStacktrace(zap.InfoLevel))
	}
	for _, writer := range options.WriterHandlers {
		zlog.WithOptions(zap.Hooks(func(e zapcore.Entry) error {
			_, err := writer.Write([]byte(e.Message))
			return err
		}))
	}
	lgr := zapr.NewLogger(zlog)
	return &logger{
		lgr: lgr,
	}
}

func (l *logger) FD() io.WriteCloser {
	return l.fd
}

func (l *logger) Fork(keysAndValues ...interface{}) Logger {
	return &logger{
		lgr: l.lgr.WithValues(keysAndValues...),
		fd:  l.FD(),
	}
}
func (l *logger) Info(msg string, ctx ...interface{}) {
	l.lgr.Info(msg, ctx...)
}
func (l *logger) V(level int) Logger {
	return &logger{
		lgr: l.lgr.V(level),
		fd:  l.FD(),
	}
}
func (l *logger) DF() io.WriteCloser {
	return l.fd
}
