package logger

import (
	"io"
	"os"

	log "github.com/inconshreveable/log15"
)

type (
	Logger interface {
		log.Logger
	}

	Options struct {
		LogToStdOut bool
		FilePath    string
	}
)

// New creates new logger based on logger.Options
func New(opt *Options) Logger {
	l := log.New()
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
		}
	}

	l.SetHandler(log.MultiHandler(handlers...))
	return l
}

// NewFromFilePath creates new logger that writes to given file
func NewFromFilePath(path string) Logger {
	return New(&Options{
		FilePath: path,
	})
}

// NewWriter creates no io.WriteCloser , it is caller responsibility to close the writer
func NewWriter(path string) (io.WriteCloser, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}
