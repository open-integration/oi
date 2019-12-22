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
		FilePath string
	}
)

// New creates new logger based on logger.Options
func New(opt *Options) Logger {
	l := log.New()
	handlers := []log.Handler{}
	if opt != nil && opt.FilePath != "" {
		h, err := log.FileHandler(opt.FilePath, log.LogfmtFormat())
		if err == nil {
			handlers = append(handlers, log.LvlFilterHandler(log.LvlDebug, h))
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
	file, err := os.Open(path)
	return file, err
}
