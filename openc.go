package core

import (
	"fmt"
	"os"

	"github.com/open-integration/open-integration/core/engine"
	"github.com/open-integration/open-integration/core/modem"
	"github.com/open-integration/open-integration/pkg/downloader"
	"github.com/open-integration/open-integration/pkg/logger"
)

var exit = func(code int) { os.Exit(code) }

type (
	// EngineOptions to create new engine
	EngineOptions struct {
		Pipeline engine.Pipeline
		// LogsDirectory path where to store logs
		LogsDirectory string
		Kubeconfig    *engine.KubernetesOptions
		Logger        logger.Logger

		serviceDownloader downloader.Downloader
		modem             modem.Modem
	}
)

// NewEngine create new engine
func NewEngine(opt *EngineOptions) engine.Engine {
	e, err := engine.New(&engine.Options{
		Pipeline:      opt.Pipeline,
		Kubeconfig:    opt.Kubeconfig,
		Logger:        opt.Logger,
		LogsDirectory: opt.LogsDirectory,
	})
	dieOnError(err)
	return e

}

func dieOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		exit(1)
	}
}

// HandleEngineError prints the error in case the engine.Run was failed and exit
func HandleEngineError(err error) {
	if err != nil {
		fmt.Printf("Failed to execute pipeline\nError: %v", err)
		exit(1)
	}
}
