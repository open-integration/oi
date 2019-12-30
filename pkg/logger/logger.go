package logger

import (
	"fmt"
	"io"
	"os"

	log "github.com/inconshreveable/log15"
)

type (
	Logger interface {
		FD() io.WriteCloser
		Debug(msg string, ctx ...interface{})
		Info(msg string, ctx ...interface{})
		Warn(msg string, ctx ...interface{})
		Error(msg string, ctx ...interface{})
		Crit(msg string, ctx ...interface{})
		New(ctx ...interface{}) Logger
	}

	Options struct {
		LogToStdOut bool
		FilePath    string
	}

	logger struct {
		logger log.Logger
		fd     io.WriteCloser
	}
)

// New creates new logger based on logger.Options
func New(opt *Options) Logger {
	l := &logger{
		logger: log.New(),
		fd:     nil,
	}
	handlers := []log.Handler{}

	lvl := log.LvlDebug
	if opt != nil {
		if opt.LogToStdOut {
			stdoutHandler := log.LvlFilterHandler(lvl, log.StdoutHandler)
			handlers = append(handlers, stdoutHandler)
		}

		if opt.FilePath != "" {
			h, err := log.FileHandler(opt.FilePath, log.LogfmtFormat())
			if err == nil {
				handlers = append(handlers, log.LvlFilterHandler(lvl, h))
			}
			f, err := os.OpenFile(opt.FilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Printf("Failed to open file %s, writer is not available\n", opt.FilePath)
			}
			l.fd = f
		}
	}

	l.logger.SetHandler(log.MultiHandler(handlers...))
	return l
}

func (l *logger) FD() io.WriteCloser {
	return l.fd
}

func (l *logger) New(ctx ...interface{}) Logger {
	return &logger{
		logger: l.logger.New(ctx),
		fd:     l.FD(),
	}
}
func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.logger.Debug(msg, ctx)
}
func (l *logger) Info(msg string, ctx ...interface{}) {
	l.logger.Info(msg, ctx)
}
func (l *logger) Warn(msg string, ctx ...interface{}) {
	l.logger.Warn(msg, ctx)
}
func (l *logger) Error(msg string, ctx ...interface{}) {
	l.logger.Error(msg, ctx)
}
func (l *logger) Crit(msg string, ctx ...interface{}) {
	l.logger.Crit(msg, ctx)
}
