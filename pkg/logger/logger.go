package logger

import (
	log "github.com/inconshreveable/log15"
	"io"
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
