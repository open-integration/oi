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
		Verbose  bool
		FilePath string
		Writer   io.Writer
	}
)

func New(opt *Options) Logger {
	l := log.New()
	lvl := log.LvlDebug
	handlers := []log.Handler{}
	verboseHandler := log.LvlFilterHandler(lvl, log.StdoutHandler)
	handlers = append(handlers, verboseHandler)

	if opt != nil && opt.FilePath != "" {
		h, err := log.FileHandler(opt.FilePath, log.LogfmtFormat())
		if err == nil {
			handlers = append(handlers, h)
		}
	}

	l.SetHandler(log.MultiHandler(handlers...))
	return l
}

func NewFromFilePath(path string) Logger {
	return New(&Options{
		FilePath: path,
	})
}

func NewWriter(path string) (io.WriteCloser, error) {
	file, err := os.Open(path)
	return file, err
}
